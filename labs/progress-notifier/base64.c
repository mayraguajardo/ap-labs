#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <inttypes.h>
#include <fcntl.h>
#include <signal.h>
#include "logger.h"

#define OUTPUT_FILE_ENCODE "-encoded.txt"
#define OUTPUT_FILE_DECODE "-decoded.txt"
#define WHITESPACE 64
#define EQUALS     65
#define INVALID    66
#define SIGINT      2

static const unsigned char d[] = {
    66,66,66,66,66,66,66,66,66,66,64,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,62,66,66,66,63,52,53,
    54,55,56,57,58,59,60,61,66,66,66,65,66,66,66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,66,66,66,66,66,66,26,27,28,
    29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66
};

int encoder(char *filename);
int decoder(char *filename);
int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize);
int base64decode(char *in, size_t inLen, unsigned char *out, size_t *outLen);
void signal_handler(int signal);

int open_file(char *filename, int mode);
off_t file_size(int fd);


off_t size;
size_t progress;
int task;

int main(int argc, char **argv){
    // Place your magic here
    initLogger("stdout");
    signal(SIGINT, signal_handler);
    signal(SIGUSR1, signal_handler);

    if(argc == 3 && strcmp(argv[1], "--encode") == 0) {
        infof("Starting Encoding File: %s\n",argv[2]);
        task = 1;
        encoder(argv[2]);
    } else if(argc == 3 && strcmp(argv[1], "--decode") == 0) {
        infof("Starting Decoding File: %s\n",argv[2]);
        task = 2;
        decoder(argv[2]);
    } else {
        errorf("Usage is ./base64 --<method> <file>");
        exit(-1);
    }

    return 0;
}

int encoder(char *filename){
    int fd;
    if((fd = open_file(filename,1)) == -1){
        errorf("Error opening %s file \nEncoding abroted \n", filename);
        exit(-1);
    }
    off_t filesize = file_size(fd);
    size = filesize;
    char *file_buffer = (char *)malloc(filesize * sizeof(char));
    char *encoded_line = (char *)malloc(filesize * sizeof(char));
    read(fd,file_buffer,filesize);
    char *token = strtok(filename, ".");
    strcat(token, OUTPUT_FILE_ENCODE);
    creat(token,0666);
    int fo = open_file(token,2);
    base64encode(file_buffer,filesize,encoded_line,filesize);
    write(fo,encoded_line,filesize);
    close(fd);
    close(fo);
}

int decoder(char *filename){
    int fd;
    if((fd = open_file(filename,1)) == -1){
        errorf("Error opening %s file \nDecoding abroted \n", filename);
        exit(-1);
    }
    off_t filesize = file_size(fd);
    size = filesize;
    size_t out_size = sizeof(char) * filesize;
    char *file_buffer = (char *)malloc(filesize * sizeof(char));
    char *encoded_line = (char *)malloc(filesize * sizeof(char));
    read(fd,file_buffer,filesize);
    char *token = strtok(filename,".");
    strcat(token, OUTPUT_FILE_DECODE);
    creat(token,0666);
    int fo = open_file(token,2);
    base64decode(file_buffer,filesize,encoded_line,&out_size);
    write(fo, encoded_line,filesize);
    close(fd);
    close(fo);
}

int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize)
{
   const char base64chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
   const uint8_t *data = (const uint8_t *)data_buf;
   size_t resultIndex = 0;
   size_t x;
   uint32_t n = 0;
   int padCount = dataLength % 3;
   uint8_t n0, n1, n2, n3;

   progress = 0;

   /* increment over the length of the string, three characters at a time */
   for (x = 0; x < dataLength; x += 3)
   {
     progress++;
     //sleep(0.001);
      /* these three 8-bit (ASCII) characters become one 24-bit number */
      n = ((uint32_t)data[x]) << 16; //parenthesis needed, compiler depending on flags can do the shifting before conversion to uint32_t, resulting to 0

      if((x+1) < dataLength)
         n += ((uint32_t)data[x+1]) << 8;//parenthesis needed, compiler depending on flags can do the shifting before conversion to uint32_t, resulting to 0

      if((x+2) < dataLength)
         n += data[x+2];

      /* this 24-bit number gets separated into four 6-bit numbers */
      n0 = (uint8_t)(n >> 18) & 63;
      n1 = (uint8_t)(n >> 12) & 63;
      n2 = (uint8_t)(n >> 6) & 63;
      n3 = (uint8_t)n & 63;

      /*
       * if we have one byte available, then its encoding is spread
       * out over two characters
       */
      if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
      result[resultIndex++] = base64chars[n0];
      if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
      result[resultIndex++] = base64chars[n1];

      /*
       * if we have only two bytes available, then their encoding is
       * spread out over three chars
       */
      if((x+1) < dataLength)
      {
         if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
         result[resultIndex++] = base64chars[n2];
      }

      /*
       * if we have all three bytes available, then their encoding is spread
       * out over four characters
       */
      if((x+2) < dataLength)
      {
         if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
         result[resultIndex++] = base64chars[n3];
      }
   }

   /*
    * create and add padding that is required if we did not have a multiple of 3
    * number of characters available
    */
   if (padCount > 0)
   {
      for (; padCount < 3; padCount++)
      {
         if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
         result[resultIndex++] = '=';
      }
   }
   if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
   result[resultIndex] = 0;
   return 0;   /* indicate success */
}

int base64decode (char *in, size_t inLen, unsigned char *out, size_t *outLen) {
    char *end = in + inLen;
    char iter = 0;
    progress = 0;
    uint32_t buf = 0;
    size_t len = 0;

    while (in < end) {
      progress++;
        unsigned char c = d[*in++];

        switch (c) {
        //case WHITESPACE: continue;   /* skip whitespace */
        case INVALID:    return 1;   /* invalid input, return error */
        case EQUALS:                 /* pad character, end of data */
            in = end;
            continue;
        default:

            buf = buf << 6 | c;
            iter++; // increment the number of iteration
            /* If the buffer is full, split it into bytes */
            if (iter == 4) {
                if ((len += 3) > *outLen) return 1; /* buffer overflow */
                *(out++) = (buf >> 16) & 255;
                *(out++) = (buf >> 8) & 255;
                *(out++) = buf & 255;
                buf = 0; iter = 0;

            }
        }
    }

    if (iter == 3) {
        if ((len += 2) > *outLen) return 1; /* buffer overflow */
        *(out++) = (buf >> 10) & 255;
        *(out++) = (buf >> 2) & 255;
    }
    else if (iter == 2) {
        if (++len > *outLen) return 1; /* buffer overflow */
        *(out++) = (buf >> 4) & 255;
    }

    *outLen = len; /* modify to reflect the actual output size */
    return 0;
}

int open_file(char *filename, int mode){
    if(mode == 1){
        int fd = open(filename, O_RDONLY);
        return fd;
    } else if(mode == 2){
        int fd = open(filename, O_WRONLY);
        return fd;
    }
    return -1;
}

off_t file_size(int fd){
    off_t current_offset = lseek(fd,(size_t)0,SEEK_CUR);
    off_t size = lseek(fd,(size_t)0, SEEK_END);
    lseek(fd,current_offset,SEEK_SET);
    return size;
}

void signal_handler(int signal){
    infof("%s file\tProgress: %f \n",(task==1)?"Coding":"Decoding",100.0f * progress / size);
}