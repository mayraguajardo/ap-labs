#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <syslog.h>
#include <signal.h>
#include "logger.h"

//Attributes
#define RESET		0
#define BRIGHT 		1
#define DIM		    2
#define UNDERLINE 	3
#define BLINK		4
#define REVERSE		7
#define HIDDEN		8

//Font color
#define BLACK   0
#define RED     1
#define GREEN   2
#define YELLOW  3
#define BLUE    4
#define MAGENTA 5
#define CYAN    6
#define	WHITE	7


int isSyslog = 0;

void font_color(int attr, int fg, int bg){
    char command[13];
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);

}
int initLogger(char *logType) {
    if(strcmp("syslog", logType) == 0){
        printf("Initializing Logger on: %s\n", logType);
        isSyslog = 1;
        return 0;
    } else {
        if (strcmp( "stdout", logType) == 0 || strcmp("stdout", logType) == 0) {
            isSyslog = 0;
            printf("Initializing Logger on: stdout\n");
            return 0;
        } else {
            errorf("Not a valid argument, got '%s'", logType);
        }
    }
}

int infof(const char *format, ...) {
    font_color(UNDERLINE,GREEN,BLACK);
    va_list args;
    va_start(args, format);
    if(isSyslog == 1){
        vsyslog(LOG_INFO,format,args);
        closelog();
    } else {
        vprintf(format,args);
        printf("\n");
    }
    va_end(args);
    font_color(RESET,WHITE,BLACK);
    return 0;
}

int warnf(const char *format, ...) {
    font_color(BRIGHT,YELLOW,BLACK);
    va_list args;
    va_start(args, format);
    if(isSyslog == 1){
        vsyslog(LOG_WARNING, format, args);
        closelog();
    } else {
        vprintf(format,args);
        printf("\n");
    }
    va_end(args);
    font_color(RESET,WHITE,BLACK);
    return 0;
}

int errorf(const char *format, ...) {
    font_color(DIM, RED, BLACK);
    va_list args;
    va_start(args, format);
    if(isSyslog == 1){
        vsyslog(LOG_ERR, format, args);
        closelog();
    } else {
        vprintf(format,args);
        printf("\n");
    }
    va_end(args);
    font_color(RESET,WHITE,BLACK);
    return 0;
}
int panicf(const char *format, ...) {
    font_color(BLINK, RED, CYAN);
    va_list args;
    va_start(args, format);
    if(isSyslog == 1){
        vsyslog(LOG_CRIT, format, args);
        closelog();
    } else {
        vprintf(format,args);
        printf("\n");
    }
    va_end(args);
    font_color(RESET,WHITE,BLACK);
    return 0;
    
}