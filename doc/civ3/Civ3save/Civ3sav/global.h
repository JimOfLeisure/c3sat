#pragma once

#define MAXLINEBUF 1024
#define MAXBUF 102400
///////////////////////////////////////////////////////////
extern FILE *fp_src, *fp_out, *fp_log;

#ifdef _DEBUG
#define debug(fmt, ...) fprintf(fp_log, fmt, __VA_ARGS__)
#else
#define debug(x)
#endif

#define output(fmt, ...) fprintf(fp_out, fmt, __VA_ARGS__)
///////////////////////////////////////////////////////////
extern void report_error(enum errorcode code);
extern int load_file();
