#include "stdafx.h"
/*
http://laurikari.net/tre/
*/
///////////////////////////////////////////////////////////
static int valid_reobj = 0;
static regex_t reobj;
static regmatch_t pmatch_global[32];
static const CHAR_T *regex_data;
static const CHAR_T *regex_pattern;
static int cflags_global;
static int avoid_eflags = 0;

///////////////////////////////////////////////////////////

int wrap_regexec(const CHAR_T *data, int eflags)
{
	regex_data = data;
	return tre_regexec(&reobj, data, elementsof(pmatch_global), pmatch_global, eflags);
}

void wrap_regfree()
{
	if (valid_reobj)
	{
		tre_regfree(&reobj);
		valid_reobj = 0;
	}
}

int wrap_regcomp(const CHAR_T *re, int flags)
{
	int errcode = 0;

	wrap_regfree();
	regex_pattern = re;
	errcode = tre_regcomp(&reobj, re, flags);

	if (errcode != 0)
	{
		printf("Comp error, regex: \"%s\"\n", regex_pattern);
		printf("	expected return code %d, got %d.\n", 0, errcode);
	} else {
		valid_reobj = 1;
	}

	return errcode;
}

void wrap_regmatch(const CHAR_T *fmt, ...)
{
	va_list ap;
	size_t pmatch_len = elementsof(pmatch_global);
	regmatch_t *pmatch = pmatch_global;
	unsigned int i;
	CHAR_T str[128];
	size_t str_len;
	CHAR_T *p;
	int *d;

	va_start(ap, fmt);

	
	for (i = 1; i < pmatch_len; i++) {
		if (*fmt++ == '%') {
			switch(*fmt++) {
			case 's':
				p = va_arg(ap, CHAR_T *);
				str_len = pmatch[i].rm_eo - pmatch[i].rm_so;
				strncpy(p, regex_data + pmatch[i].rm_so, str_len);
				p[str_len] = 0;
				break;
			case 'd':
				d = va_arg(ap, int *);
				str_len = pmatch[i].rm_eo - pmatch[i].rm_so;
				strncpy(str, regex_data + pmatch[i].rm_so, str_len);
				str[str_len] = 0;
				*d = atoi(str);
				break;
			default:
				goto _end;
			}
		}
	}

_end:

	va_end(ap);
}