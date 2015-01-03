// Civ3sav.cpp : 定义控制台应用程序的入口点。
//

#include "stdafx.h"
#include "global.h"

///////////////////////////////////////////////////////////
FILE *fp_src, *fp_out, *fp_log;

///////////////////////////////////////////////////////////
void print_usage()
{
	printf("civ3save *.sav *.txt log.txt\n");
}


int main(int argc, char* argv[])
{
	if (argc != 3) {
		goto err;
	}
	fp_src = fopen(argv[1], "rb");
	if (!fp_src)
		goto err;
	fp_out = fopen(argv[2], "w");
	if (!fp_out)
		goto err;
	fp_log = fopen("log.txt", "w");
	if (!fp_log)
		goto err;

	load_file();

	fclose(fp_log);
	fclose(fp_out);
	fclose(fp_src);
	return 0;

err:
	fclose(fp_log);
	fclose(fp_out);
	fclose(fp_src);
	print_usage();
	return -1;
}


