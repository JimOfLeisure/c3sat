#include "stdafx.h"
#include "Civ3global.h"

///////////////////////////////////////////////////////////
extern char linebuf[MAXLINEBUF];

enum errorcode g_errorcode;
enum token_line looknext;

extern char section[5];

extern int	bytes;
extern char type[MAXLINEBUF];
extern char field[MAXLINEBUF];

extern char g_buf[MAXLINEBUF];

char g_buf_c[MAXLINEBUF];
char *cache1;
char *cache2;
char *cache3;
int cache1_off;
int cache2_off;
int cache3_off;

int field1;
int field2;
int field3;

enum cache {
	CACHE1 = 1,
	CACHE2,
	CACHE3,
};
///////////////////////////////////////////////////////////
extern enum token_line parse_line();

void report_error(enum errorcode code)
{
	printf("errorcode = %d\n", (int)code);
	printf("%s\n", linebuf);
    g_errorcode = code;
}

int sizeof_type(char *t)
{
	if (!strcmp(t, "char")) {
		return 1;
	} else if (!strcmp(t, "long")) {
		return 4;
	} else if (!strcmp(t, "string")) {
		return 1;
	} else if (!strcmp(t, "byte")) {
		return 1;
	} else {
		return 1;
	};
}

char* to_c(char *token, char *c_token, int toupper)
{
	char *c_t = c_token;
	char *t = token;
	if (isdigit(*t)) {
		*c_t++ = '_';
	}
	for (; *t; ++t, ++c_t) {
		if (!isalnum(*t)) {
			*c_t = '_';
		} else {
			*c_t = toupper ? _toupper(*t) : *t;
		}
	}
	*c_t = 0;

	return c_token;
}

///////////////////////////////////////////////////////////
void get_cache (enum cache cache_t, char **cache, int **cache_off)
{
	switch (cache_t) {
	case CACHE1:
		*cache = cache1;
		*cache_off = &cache1_off;
		break;
	case CACHE2:
		*cache = cache2;
		*cache_off = &cache2_off;
		break;
	case CACHE3:
		*cache = cache3;
		*cache_off = &cache3_off;
		break;
	}
}
void rewind_cache (enum cache cache_t)
{
	char *cache;
	int *cache_off;
	get_cache(cache_t, &cache, &cache_off);
	*cache_off = 0;
}
int output_cache (enum cache cache_t, char* buf, int len)
{
	char *cache;
	int *cache_off;
	get_cache(cache_t, &cache, &cache_off);

	memcpy(cache + *cache_off, buf, len);
	*cache_off += len;
	cache[*cache_off] = 0;

	return len;
}
void flush_cache (enum cache cache_t)
{
	switch (cache_t) {
	case CACHE1:
		fwrite(cache1, cache1_off, 1, fp_out);
		break;
	case CACHE2:
		output_cache(CACHE1, cache2, cache2_off);
		cache2_off = 0;
		break;
	case CACHE3:
		output_cache(CACHE2, cache3, cache3_off);
		cache3_off = 0;
		break;
	}
}
///////////////////////////////////////////////////////////

void export_line(enum token_line t)
{
	int count;
	char buf[MAXLINEBUF];
	char buf1[MAXLINEBUF];
	char cache[MAXLINEBUF];
	int len;

	switch (t) {
	case SECTION_HEADER:
		output("struct sect%s {\n", to_c(section, buf1, 0));
		break;
	case SECTION_NAME:
		output("\t%s\tname[%d];\n", to_c(type, buf1, 0), bytes);
		break;
	case SECTION_NUMBER:
		count = bytes / sizeof_type(type);
		if (count == 1) {
			output("\t%s\tnumber;\n", to_c(type, buf1, 0));
		} else {
			output("\t%s\tnumber[%d];\n", to_c(type, buf1, 0), count);
		}
		break;
	case ITEM_HEADER:
		output("\tstruct st%s *%s;\n", to_c(g_buf, buf, 1), to_c(g_buf, buf1, 0));
		rewind_cache(CACHE1);
		field1 = 0;

		len = sprintf(cache, "\tstruct st%s {\n", to_c(g_buf, buf, 1));
		output_cache(CACHE1, cache, len);
		break;
	case ITEM_FIELD:
		count = bytes / sizeof_type(type);
		if (count == 1) {
			len = sprintf(cache, "\t\t%s\tf%d_%s;\n", to_c(type, buf, 0), field1++, to_c(field, buf1, 0));
		} else {
			len = sprintf(cache, "\t\t%s\tf%d_%s[%d];\n", to_c(type, buf, 0), field1++, to_c(field, buf1, 0), count);
		}
		output_cache(CACHE1, cache, len);
		break;
	case ITEM_FIELD_BINARY:
		break;
	case ITEM_FIELD_BINARY_HEADER:
		break;
	case ITEM_FIELD_BINARY_BITFIELD:
		break;
	case ITEM_ITEM_HEADER:
		len = sprintf(cache, "\t\tstruct stt%s *%s;\n", to_c(g_buf, buf, 1), to_c(g_buf, buf1, 0));
		output_cache(CACHE1, cache, len);
		rewind_cache(CACHE2);
		field2 = 0;

		len = sprintf(cache, "\t\tstruct stt%s {\n", to_c(g_buf, buf, 1));
		output_cache(CACHE2, cache, len);
		break;
	case ITEM_ITEM_FIELD:
		count = bytes / sizeof_type(type);
		if (count == 1) {
			len = sprintf(cache, "\t\t\t%s\tf%d_%s;\n", to_c(type, buf, 0), field2++, to_c(field, buf1, 0));
		} else {
			len = sprintf(cache, "\t\t\t%s\tf%d_%s[%d];\n", to_c(type, buf, 0), field2++, to_c(field, buf1, 0), count);
		}
		output_cache(CACHE2, cache, len);
		break;
	case ITEM_ITEM_ITEM_HEADER:
		len = sprintf(cache, "\t\t\tstruct sttt%s *%s;\n", to_c(g_buf, buf, 1), to_c(g_buf, buf1, 0));
		output_cache(CACHE2, cache, len);
		rewind_cache(CACHE3);
		field3 = 0;

		len = sprintf(cache, "\t\t\tstruct sttt%s {\n", to_c(g_buf, buf, 1));
		output_cache(CACHE3, cache, len);
		break;
	case ITEM_ITEM_ITEM_FIELD:
		count = bytes / sizeof_type(type);
		if (count == 1) {
			len = sprintf(cache, "\t\t\t\t%s\tf%d_%s;\n", to_c(type, buf, 0), field3++, to_c(field, buf1, 0));
		} else {
			len = sprintf(cache, "\t\t\t\t%s\tf%d_%s[%d];\n", to_c(type, buf, 0), field3++, to_c(field, buf1, 0), count);
		}
		output_cache(CACHE3, cache, len);
		break;
	default:
		;
	}
}

void match_line(enum token_line t)
{
	if (looknext == t) {
		export_line(t);
        looknext = parse_line();
	}
    else
        report_error(SYNTAX);
}

int parse_item_item_item()
{
	match_line(ITEM_ITEM_ITEM_HEADER);
	match_line(ITEM_ITEM_ITEM_FIELD);
	while(1) {
		switch (looknext) {
		case ITEM_ITEM_ITEM_FIELD:
			match_line(ITEM_ITEM_ITEM_FIELD);
			break;
		default:
			goto _end;
		}
	}
_end:
	return 0;
}

int parse_item_item()
{
	match_line(ITEM_ITEM_HEADER);
	match_line(ITEM_ITEM_FIELD);
	while(1) {
		switch (looknext) {
		case ITEM_ITEM_FIELD:
			match_line(ITEM_ITEM_FIELD);
			break;
		case ITEM_ITEM_ITEM_HEADER:
			parse_item_item_item();
		default:
			goto _end;
		}
	}
_end:
	output_cache(CACHE2, "\t\t};\n", 5);
	flush_cache(CACHE3);
	return 0;
}

int parse_item()
{
	match_line(ITEM_HEADER);
	match_line(ITEM_FIELD);
	while(1) {
		switch (looknext) {
		case ITEM_FIELD:
			match_line(ITEM_FIELD);
			break;
		case ITEM_ITEM_HEADER:
			parse_item_item();
			break;
		default:
			goto _end;
		}
	}
_end:
	output_cache(CACHE1, "\t};\n", 4);
	flush_cache(CACHE2);
	return 0;
}

int parse_section()
{
	match_line(SECTION_HEADER);
	match_line(SECTION_NAME);
	match_line(SECTION_NUMBER);

	parse_item();

	output("};\n");
	flush_cache(CACHE1);
	return 0;
}


int parse_file()
{
	cache1 = (char *)malloc(MAXBUF);
	cache2 = (char *)malloc(MAXBUF);
	cache3 = (char *)malloc(MAXBUF);

	output("typedef char * string;\n");
	output("typedef char byte;\n");

	while (!feof(fp_src)) {

		parse_line();
		switch (looknext) {
		case SECTION_HEADER:
			parse_section();
			break;
//		default:
//			parse_line();
		}



//		fputs(linebuf, fp_out);
	}
	free(cache1);
	free(cache2);
	free(cache3);
	cache1 = NULL;
	cache2 = NULL;
	cache3 = NULL;
	return 0;
}

///////////////////////////////////////////////////////////