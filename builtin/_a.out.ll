;package main

declare i32 @tiny_go_builtin_exit(i32)
declare i32 @tiny_go_builtin_println(i32)

define i32 @tiny_go_main_init() {
	ret i32 0
}
define i32 @tiny_go_main_main() {
	%t0 = add i8 0, 97
	%local_a.pos.56 = alloca i8, align 1
	store i8 %t0, i8* %local_a.pos.56
	%t1 = add i8 0, 98
	%local_b.pos.68 = alloca i8, align 1
	store i8 %t1, i8* %local_b.pos.68
	%t3 = load i8, i8* %local_a.pos.56, align 1
	%t4 = load i8, i8* %local_b.pos.68, align 1
	%t5 = sext i8 %t3 to i32
	%t6 = sext i8 %t4 to i32
	%t2 = add i32 %t5, %t6
	%t7 = trunc i32 %t2 to i8
	%local_c.pos.80 = alloca i8, align 1
	store i8 %t7, i8* %local_c.pos.80
	ret i32 0
}

define i32 @main() {
	call i32() @tiny_go_main_init()
	call i32() @tiny_go_main_main()
	ret i32 0
}
