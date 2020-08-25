#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* print Fahrenheit-Celsius table */
int main(int argc, char **argv)
{
   int fahr;
    if (argc-1 == 1){
        fahr = atoi(argv[1]);
        printf("Fahrenheit: %3d, Celsius: %6.1f\n",fahr,(5.0/9.0)*(fahr-32));
    }
    else if (argc-1 == 3){
        int LOWER = atoi(argv[1]);
        int UPPER = atoi(argv[2]);
        int STEP = atoi(argv[3]);

        for (fahr = LOWER; fahr <= UPPER; fahr = fahr + STEP)
	    printf("Fahrenheit: %3d, Celcius: %6.1f\n", fahr, (5.0/9.0)*(fahr-32));
    } 
    else{
        printf("No valid arguments\n");
    }

    return 0;
}
