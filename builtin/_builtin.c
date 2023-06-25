#include <stdio.h>
#include <stdlib.h>

int tiny_go_builtin_println(int x){
    return printf("%d\n",x);
}

int tiny_go_builtin_exit(int x){
    exit(x);
    return 0;
}