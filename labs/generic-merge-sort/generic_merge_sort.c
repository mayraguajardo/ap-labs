#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define MAX_LINES 10000

void merge_sort(void *array[], int left, int right, int (*comp)(void *, void*));
void merge(void *array[], int left, int middle, int right, int (*comp)(void *, void *));
int numcmp(const char *s1, const char *s2);
void print_merged(int size, char *data[]);

char *lines[MAX_LINES];

int main(int argc, char** argv)
{
    int numeric = 0;
    FILE *file;

    if(argc < 2) {
        printf("Usage: ./generic_merge_sort (-n) <filename>\n");
        return 0;
    }else if(argc == 2 && strcmp(argv[1],"-n") == 0){
        printf("File name missing \n");
        return 0;
    }

    if(argc == 3 && strcmp(argv[1],"-n") == 0)
        numeric = 1;
    
    if((file = fopen(argv[2],"r")) == NULL && (file = fopen(argv[1], "r")) == NULL){
        perror("An error ocurred while opening file\n");
        return 0;
    }

    char line[256] = {0};
    int n = 0;

    while(fgets(line,100,file)){
        lines[n] = (char*)malloc(strlen(line) + sizeof(char*));
        strcpy(lines[n], line);
        n++;
    }

    if(fclose(file))
        perror("An error ocurred while closing file");
    
    merge_sort((void *) lines, 0,n-1, (int (*)(void*, void*)) (numeric ? numcmp : strcmp));
    print_merged(n, lines);
    return 0;
}

void merge_sort(void *array[], int left, int right, int(*comp)(void *, void *)){
    if(left < right){
        int middle = (left + right) / 2;

        merge_sort(array, left, middle, (*comp));
        merge_sort(array, middle + 1, right, (*comp));
        merge(array, left, middle, right, (*comp));

    }
}

void merge(void *array[], int left, int middle, int right, int (*comp)(void *, void *)){
    int left_side = middle - left + 1;
    int right_side = right - middle;

    char *left_array[left_side], *right_array[right_side];
    for (int i = 0; i < left_side; i++)
        left_array[i] = array[left + i];

    for (int i = 0; i < right_side; i++)
        right_array[i] = array[middle + i + 1];
    
    int i = 0;
    int j = 0;
    int k = left;
    while (i < left_side && j < right_side){
        if((*comp) (left_array[i], right_array[j]) < 0){
            array[k] = left_array[i];
            i++;
        } else {
            array[k] = left_array[j];
            j++;
        }
        k++;
    }

    while (i < left_side){
        array[k] = left_array[i];
        i++;
        k++;
    }
    while(j < right_side){
        array[k] = right_array[j];
        j++;
        k++;
    }
}
int numcmp (const char *s1, const char *s2){
        double v1 = atof(s1);
        double v2 = atof(s2);
        if(v1 < v2)
            return -1;
        else if (v1 > v2)
            return 1;
        else
            return 0;
    }
void print_merged(int size, char *data[]){
    for(int i = 0; i < size; i++)
        printf("%s", data[i]);
    printf("\n");
}