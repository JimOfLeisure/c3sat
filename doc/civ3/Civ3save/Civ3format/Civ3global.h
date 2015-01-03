#pragma once

enum token_line {
	UNKNOWN_TOKEN = 0,
	SECTION_HEADER,
	SECTION_NAME,
	SECTION_NUMBER,
	ITEM_HEADER,
	ITEM_FIELD,
	ITEM_FIELD_BINARY,
	ITEM_FIELD_BINARY_HEADER,
	ITEM_FIELD_BINARY_BITFIELD,
	ITEM_ITEM_HEADER,
	ITEM_ITEM_FIELD,
	ITEM_ITEM_ITEM_HEADER,
	ITEM_ITEM_ITEM_FIELD,
	COMMENT,
	OTHERS,
};
enum errorcode {
	SYNTAX = 0,
};
typedef enum token_line (parse_line_func_t)();

#ifdef _DEBUG
#define debug(fmt, ...) fprintf(fp_log, fmt, __VA_ARGS__)
#else
#define debug(x)
#endif

#define output(fmt, ...) fprintf(fp_out, fmt, __VA_ARGS__)

#define MAXLINEBUF 1024
#define MAXBUF 102400
///////////////////////////////////////////////////////////
extern FILE *fp_src, *fp_out, *fp_log;

extern enum errorcode g_errorcode;
extern enum token_line looknext;

///////////////////////////////////////////////////////////
extern void report_error(enum errorcode code);
extern int parse_file();

extern int wrap_regcomp(const CHAR_T *re, int flags);
extern int wrap_regexec(const CHAR_T *data, int eflags);
extern void wrap_regfree();
extern void wrap_regmatch(const CHAR_T *fmt, ...);
