#include <stdio.h>
#include <stdlib.h>
static char daytab[2][13]={
    {0,31,28,31,30,31,30,31,31,30,31,30,31},
    {0,31,29,31,30,31,30,31,31,30,31,30,31}
};

int day_of_the_year(int year, int month, int day){ 
    int i, leap;
    leap = year%4 == 0 && year%100 != 0 || year%400 == 0;

    for (i =1; i<month; i++)
        day += daytab[leap][i];
    return day; 
}


/* month_day function's prototype*/
void month_day(int year, int yearday, int *pmonth, int *pday){

    int i, leap;
    if(year <=0){
        *pmonth = 0;
        *pday = 0;
        return;
    }
    leap = year%4 == 0 && year%100 != 0 || year%400 == 0;



    for(i=1;i<=12 && yearday > daytab[leap][i];i++)
        yearday -= daytab[leap][i];

    if(i > 12 && yearday > daytab[leap][12])
    {
        *pmonth=0;
        *pday=0;
    }
    else
    {
        *pmonth=i;
        *pday=yearday;
    }

    


}

int main(int argc, char **argv) {

    int year = atoi(argv[1]);
    int yearday = atoi(argv[2]);

    int day,mon;
    month_day(year,yearday,&mon,&day);
    printf("Month: %d, Day: %d\n", mon,day);


    return 0;
}
