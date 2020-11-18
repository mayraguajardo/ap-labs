#include <stdio.h>
#include <math.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>
#include <pthread.h>
#include <sys/stat.h>
#include <sys/types.h>
#include "logger.h"

#define MAT_DIM 2000
#define NUM_THREADS 2000

int NUM_BUFFERS;
char *RESULT_MATRIX_FILE;
long *result;
long **buffers;

pthread_mutex_t *mutex;
pthread_t threads[NUM_THREADS];

typedef struct {
    int row;
    int col;
    long *matA;
    long *matB;
}thread_args;

long *readMatrix(char *filename);
long *getColumn(int col, long *matrix);
long *getRow(int row, long *matrix);
int getLock();
int releaseLock(int lock);
long dotProduct(long *vect1, long *vect2);
long *multiply(long *matA, long *matB);
int saveResultMatrix(long *result);


int main(int argc, char** argv){
    // Place your magic here

    if(argc != 5){
        errorf("Wrong arguments");
        exit(1);

    }else{
        for(int i = 1; i < argc; i ++){
            if(strcmp(argv[i],"-n") == 0){
                i++;
                if(i >= argc){
                    errorf("Wrong arguments");
                    exit(1);
                }
                NUM_BUFFERS = atoi(argv[i]);
            } else if(strcmp(argv[i], "-out") == 0){
                i++;
                if(i>=argc){
                    errorf("Wrong arguments");
                    exit(1);
                }
                RESULT_MATRIX_FILE = argv[i];
            }

        }

    }

    buffers = malloc(NUM_BUFFERS*sizeof(long *));
    mutex = malloc(NUM_BUFFERS * sizeof(pthread_mutex_t));
    long *matrixA = readMatrix("matA.dat");
    infof("Matrix A read\n");
    long *matrixB = readMatrix("matB.dat");
    infof("Matrix B read\n");
    for(int i = 0; i < NUM_BUFFERS; i++){
        pthread_mutex_init(&mutex[i],NULL);
    }

    long *result = multiply(matrixA, matrixB);
    infof("Multiplication done\n");
    infof("Saving result\n");
    int savedResult = saveResultMatrix(result);
    free(matrixA);
    free(matrixB);
    free(mutex);
    free(buffers);

    if (savedResult != 0){
        panicf("There was an error saving the file");
    }
    return 0;

}

long multiplication(thread_args *arg){
    int buffer1 = -1,
        buffer2 = -1;
    while(buffer1 == -1 || buffer2 == -1){
        if(buffer1 == -1){
            buffer1 = getLock();
        }
        if(buffer2 == -1){
            buffer2 = getLock();
        }
        
    }
    buffers[buffer1] = getRow(arg->row,arg->matA);
    buffers[buffer2] = getColumn(arg->col,arg->matB);
    long result = dotProduct(buffers[buffer1], buffers[buffer2]);
    free(buffers[buffer1]);
    free(buffers[buffer2]);
    free(arg);
    releaseLock(buffer1);
    releaseLock(buffer2);
    return result;
}

long *multiply(long *matrixA, long *matrixB){
    infof("Starting multiplication. It may take a while, be patient :)");
    long *result = malloc(MAT_DIM*MAT_DIM*sizeof(long));
    for(int i = 0; i < MAT_DIM; i++){
        for(int j = 0; j < MAT_DIM; j++){
            thread_args *ar = malloc(sizeof(thread_args));
            ar -> row = i;
            ar -> col = j;
            ar -> matA = matrixA;
            ar -> matB = matrixB;
            pthread_create(&threads[j],NULL,(void * (*)(void *))multiplication,(void *)ar);
        }
        for(int j = 0; j < MAT_DIM; j++){
            void *res;
            pthread_join(threads[j], &res);
            result[MAT_DIM * j + i] = (long) res;
        }
    }
    return result;
}

long * readMatrix(char *filename){
    int size = 0;
    FILE *f = fopen(filename, "r");
    if(f == NULL){
        errorf("Error opening the file");
        exit(2);
    }
    char c;
    while((c = fgetc(f)) != EOF){
        if(c == '\n')
            size++;
    }
    rewind(f);
    long *matrix = malloc(size*sizeof(long));
    int index = 0;
    while(fscanf(f, "%ld", &matrix[index]) != EOF){
        index++;
    }
    fclose(f);
    return matrix;
}

long *getRow(int row, long *matrix){
    long *values = (long *)malloc(MAT_DIM*sizeof(long));
    long actualPos = 0;
    for(int i = MAT_DIM*row; i < MAT_DIM*row+MAT_DIM; i++)
        values[actualPos++] = matrix[i];
    return values;
}

long *getColumn(int col, long *matrix){
    long *values = (long *)malloc(MAT_DIM*sizeof(long));
    long actualPos = 0;
    for(int i = col; i < MAT_DIM*MAT_DIM; i += MAT_DIM)
        values[actualPos++] = matrix[i];
    return values;
}

int getLock(){
    for(int i = 0; i < NUM_BUFFERS; i++){
        if(pthread_mutex_trylock(&mutex[i]) == 0){
            return i;
        }
    }
    return -1;
}

int releaseLock(int lock){
    return pthread_mutex_unlock(&mutex[lock]);
}

long dotProduct(long *vect1, long *vect2){
    long result = 0;
    for(int i = 0; i < MAT_DIM; i++)
        result += vect1[i] * vect2[i];
    return result;
}

int saveResultMatrix(long *result){
    FILE *file;
    file = fopen(RESULT_MATRIX_FILE, "w+");
    if(file == NULL)
        return -1;
    int size = MAT_DIM * MAT_DIM;
    for(int i = 0; i < size; i++)
        fprintf(file, "%ld\n", result[i]);
    fclose(file);
    return 0;
}