#include <stdio.h>

void reverse(void)
{
    char c;
    if((c = getchar()) != '\n'){ reverse(); }
    putchar(c);
    return;
}
int main(){
    // Place your magic here

    printf("Enter a string:\n");
    reverse();
    putchar('\n'); 
    return 0;

}
