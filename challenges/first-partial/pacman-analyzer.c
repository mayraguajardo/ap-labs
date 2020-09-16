#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <sys/stat.h>


struct package {
    char name[200];
    char install_date[200];
    char last_update[200];
    int updates;
    char removal_date[200];
};

struct package packages[5000];
void analizeLog(char *logFile, char *report);
int get_line(FILE *file, char *buffer, size_t size);
int package_type(char* line);
char* get_name(char* line);
char* get_date(char* line);

int main(int argc, char **argv) {

    if (argc < 5 || strcmp(argv[1],"-input") != 0 || strcmp(argv[3],"-report") != 0) {
        printf("Usage:./pacman-analizer.o -input <input_file.txt> -report <report_file.txt>\n");
        return 1;
    }
    analizeLog(argv[2],argv[4]);

    return 0;
}

int get_line(FILE *file, char *buffer, size_t size){
    if (size == 0)
        return 0;
    
    size_t current_size = 0;
    int c;

    while((c = (char) getc(file)) != '\n' && current_size < size){
        if(c == EOF){
            break;
        }
        buffer[current_size] = (char) c;
        current_size++;
    }
    if (current_size == 0)
        return 0;
    
    buffer[current_size] = '\0';
    return current_size;
}

char* get_date(char* line){
    int size = 0;
    for (size; line[size] != ']'; size++);
    size++;
    char *date = calloc(1,size);
    int cont = 0;
    for (int i =0; i < size; i++, cont++){
        date[cont] = line[i];
    }
    return date;
}

char* get_name(char* line){
    int cont = 0, start_point = 0, size = 0;

    for(int i = 0; i<2; i++){
        for (start_point; line[start_point] != ']'; start_point++);
        start_point += 2;
    }

    for (start_point; line[start_point] != ' '; start_point++);
    start_point++;
    for(int j = start_point + 1; line[j] != ' '; j++){
        size++;
    }
    char *name = calloc(1,size);

    for(int k = start_point; line[k] != ' '; k++, cont++){
        name[cont] = line[k];
    }
    return name;
    
}

int package_type (char* line){
    int start_point = 0;
    for( int i = 0; i < 2; i++){
        for(start_point; line[start_point] != '\0'; start_point++){
            if(line[start_point] == ']')
                break;
            
        }
        start_point += 2;
    }

    if(line[start_point] == 'i' && line[start_point + 1] == 'n' && line[start_point + 2] == 's' )
        return 1;
    if(line[start_point] == 'u' && line[start_point + 1] == 'p' && line[start_point + 2] == 'g' )
        return 2;
    if(line[start_point] == 'r' && line[start_point + 1] == 'e' && line[start_point + 2] == 'm' )
        return 3;
    return 0;
}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] log file\n", logFile);

    // Implement your solution here.

    char line[255];
    int c;

    FILE* file;
    file = fopen(logFile, "r");

    if(file == NULL){
        perror("Error opening the log file\n");
        exit(EXIT_FAILURE);
    }

    
    int writer = open(report, O_WRONLY|O_CREAT|O_TRUNC,0644);
    if (writer < 0){
        perror("Error opening/creating the report file");
        exit(EXIT_FAILURE);
    }

    int installed = 0, 
        removed = 0, 
        upgraded = 0, 
        current = 0;
    while (c = get_line(file,line,255) > 0){
        int n = package_type(line);
        if(n==1){
            char* name = get_name(line);
            char* date = get_date(line);
            strcpy(packages[current].name,name);
            strcpy(packages[current].install_date,date);
            packages[current].updates = 0;
            strcpy(packages[current].removal_date, "-");
            current++;
            installed++;
        } else if (n==2){
            char* name = get_name(line);
            char* date = get_date(line);
            for (int i =0; i < 1500; i++){
                if(strcmp(packages[i].name,name) == 0){
                    strcpy(packages[i].last_update,date);
                    if(packages[i].updates == 0)
                        upgraded++;
                    packages[i].updates++;
                    break;
                }
            }
        } else if(n==3){
            char* name = get_name(line);
            char* date = get_date(line);
            for(int i = 0; i<1500; i++){
                if(strcmp(packages[i].name,name) == 0)
                    strcpy(packages[i].removal_date,date);
                break;
            }
            removed++;
        }
        
    }

    write(writer, "Pacman Packages Report\n", strlen("Pacman Packages Report\n"));
    write(writer,"----------------------\n",strlen("----------------------\n"));
    char aux[10];
    write(writer, "Installed packages : ", strlen("Installed packages : "));
    sprintf(aux, "%d\n", installed);
    write(writer, aux, strlen(aux));
    write(writer, "Upgraded packages : ",strlen("Upgraded packages : "));
    sprintf(aux, "%d\n", upgraded);
    write(writer, aux, strlen(aux));
    write(writer, "Removed packages : ",strlen("Removed packages : "));
    sprintf(aux, "%d\n", removed);
    write(writer, aux, strlen(aux));
    write(writer, "Current installed : \n",strlen("Current installed : "));
    sprintf(aux, "%d\n", (installed-removed));
    write(writer, aux, strlen(aux));

    write(writer, "\n\nList of packages\n", strlen("\n\nList of packages\n"));
    write(writer,"----------------------\n",strlen("----------------------\n"));
    for(int i = 0; i < 1500; i++){
        if(strcmp(packages[i].name, "") != 0){
            write(writer, "- Package name         : ",strlen("- Package name         : "));
            write(writer,packages[i].name, strlen(packages[i].name));
            write(writer, "\n   - Install date      : ",strlen("\n   - Install date      : "));
            write(writer,packages[i].install_date, strlen(packages[i].install_date));
            write(writer, "\n   - Last update date  : ",strlen("\n   - Last update date  : "));
            write(writer,packages[i].last_update, strlen(packages[i].last_update));
            write(writer, "\n   - How many updates  : ",strlen("\n   - How many updates  : "));
            sprintf(aux, "%d", packages[i].updates);
            write(writer,aux, strlen(aux));
            write(writer, "\n   - Removal date      : ",strlen("\n   - Removal date      : "));
            write(writer,packages[i].removal_date, strlen(packages[i].removal_date));
            write(writer, "\n",strlen("\n"));
        } else if (strcmp(packages[i].name, "") == 0){
            break;
        }
    }

    if(close(writer) < 0){
        perror("Error trying to close file");
        exit(1);
    }



    printf("Report is generated at: [%s]\n", report);
}
