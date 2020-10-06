#include <sys/inotify.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <stdint.h>
#include <ftw.h>
#include <fcntl.h> 
#include "logger.h"

#define _XOPEN_SOURCE 500
#define BUFFER_LEN  (10 * (sizeof(struct inotify_event) + 1024 + 1))
#define MAX_EVENTS 250

int inotify_fd, wd;
char buf[BUFFER_LEN] __attribute__ ((aligned(8)));
ssize_t num_read;
char *p;    
struct inotify_event * event;
int new_directory = 0;

static void display_Inotify_event(struct inotify_event * i){
    char old_name[1024];
    if( i-> mask & IN_CREATE){
        if(event -> mask & IN_ISDIR){
            infof("- [Directory - Create] - %s\n", i -> name);
            new_directory = 1;
        } else
            infof("- [File - Create] - %s\n",i->name);
        
    }

    if(i -> mask & IN_DELETE){
        if(event -> mask & IN_ISDIR)
            infof("- [Directory - Removal] - %s\n",i->name);
        else
            infof("- [File - Removal] - %s\n",i->name);
        
    }

    if( i -> mask & IN_MODIFY){
        if(event -> mask & IN_ISDIR)
            infof("- [Directory - Modify] - %s\n",i->name);
        else
            infof("- [File - Modify] - %s\n",i->name);
        
    }

    if(i -> mask & IN_MOVED_FROM)
        strcpy(old_name, i->name);
    
    if(i -> mask & IN_MOVED_TO){
        if(i -> cookie > 0){
            if(event -> mask & IN_ISDIR)
                infof("- [Directory - Rename] - %s\n",i->name);
            else
                infof("- [File - Rename] - %s\n",i->name);
        }
        memset(old_name, 0, strlen(old_name));
    }
}

int watch_documents(const char *path_name, const struct stat *sub_dirs, int tflag){
    wd = inotify_add_watch(inotify_fd,path_name,IN_ALL_EVENTS);
    if(wd == -1)
        return -1;
    if(new_directory == 0)
        infof("Watching : %s\n"),path_name;
    return 0;
}




int main(int argc, char *argv[]){
    if(argc < 2 || strcmp(argv[1], "--help") == 0)
        errorf("Usage is: %s <PATHNAME>\n", argv[0]);
    if (argc < 2)
        exit(-1);
    if(argc == 2){
        inotify_fd = inotify_init();
        if(inotify_fd == -1){
            errorf("Error initializing inotify");
            exit(0);
        }
        ftw(argv[1], watch_documents, 2048);
        for(;;){
            num_read = read(inotify_fd, buf, BUFFER_LEN);
            if(num_read == 0)
                infof("read() from inotify fd returned 0");

            for(p = buf; p < buf + num_read; ){
                event = (struct inotify_event * ) p;
                display_Inotify_event(event);
                p += sizeof(struct inotify_event) + event -> len;
            }

            if(new_directory == 1){
                ftw(argv[1], watch_documents, 2048);
                new_directory = 0;
            }    
        }
        inotify_rm_watch(inotify_fd, wd);
        close(inotify_fd);
        return 0;
    }

    
    return 0;
}
