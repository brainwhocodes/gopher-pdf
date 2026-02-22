//go:build cgo && !nocgo

package fitz

/*
#include <mupdf/fitz.h>
#include <stdlib.h>
#include <string.h>

const char *fz_version = FZ_VERSION;

typedef struct pdf_document pdf_document;
typedef struct pdf_obj pdf_obj;

pdf_obj *pdf_load_object(fz_context *ctx, pdf_document *doc, int num);
void pdf_print_obj(fz_context *ctx, fz_output *out, pdf_obj *obj, int tight, int ascii);
void pdf_update_object(fz_context *ctx, pdf_document *doc, int num, pdf_obj *newobj);
void pdf_drop_obj(fz_context *ctx, pdf_obj *obj);
void pdf_dict_putp(fz_context *ctx, pdf_obj *dict, const char *path, pdf_obj *val);
void pdf_dict_dels(fz_context *ctx, pdf_obj *dict, const char *key);
pdf_obj *pdf_dict_getp(fz_context *ctx, pdf_obj *dict, const char *path);
int pdf_is_dict(fz_context *ctx, pdf_obj *obj);

fz_buffer *pdf_load_stream_number(fz_context *ctx, pdf_document *doc, int num);
fz_buffer *pdf_load_raw_stream_number(fz_context *ctx, pdf_document *doc, int num);

enum
{
	PDF_LEXBUF_SMALL = 256
};

typedef struct
{
	size_t size;
	size_t base_size;
	size_t len;
	int64_t i;
	float f;
	char *scratch;
	char buffer[PDF_LEXBUF_SMALL];
} pdf_lexbuf;

void pdf_lexbuf_init(fz_context *ctx, pdf_lexbuf *lexbuf, int size);
void pdf_lexbuf_fin(fz_context *ctx, pdf_lexbuf *lexbuf);
pdf_obj *pdf_parse_stm_obj(fz_context *ctx, pdf_document *doc, fz_stream *f, pdf_lexbuf *buf);

pdf_obj *safe_pdf_load_object(fz_context *ctx, pdf_document *pdf, int num)
{
	pdf_obj *obj;

	fz_try(ctx)
	{
		obj = pdf_load_object(ctx, pdf, num);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return obj;
}

int safe_pdf_print_obj(fz_context *ctx, fz_output *out, pdf_obj *obj)
{
	fz_try(ctx)
	{
		pdf_print_obj(ctx, out, obj, 1, 0);
	}
	fz_catch(ctx)
	{
		return 0;
	}

	return 1;
}

static pdf_obj *safe_pdf_parse_stm_obj_from_string(fz_context *ctx, pdf_document *pdf, const char *text)
{
	fz_stream *stm = NULL;
	pdf_obj *obj = NULL;
	pdf_lexbuf lexbuf;
	int lexbuf_init = 0;

	if (!text)
		return NULL;

	fz_try(ctx)
	{
		stm = fz_open_memory(ctx, (const unsigned char *)text, strlen(text));
		pdf_lexbuf_init(ctx, &lexbuf, PDF_LEXBUF_SMALL);
		lexbuf_init = 1;
		obj = pdf_parse_stm_obj(ctx, pdf, stm, &lexbuf);
	}
	fz_always(ctx)
	{
		if (lexbuf_init)
			pdf_lexbuf_fin(ctx, &lexbuf);
		if (stm)
			fz_drop_stream(ctx, stm);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return obj;
}

int safe_pdf_dict_putp_parsed(fz_context *ctx, pdf_document *pdf, pdf_obj *dict, const char *path, const char *value)
{
	pdf_obj *val = NULL;

	fz_try(ctx)
	{
		val = safe_pdf_parse_stm_obj_from_string(ctx, pdf, value);
		if (!val)
			fz_throw(ctx, FZ_ERROR_ARGUMENT, "bad value");
		pdf_dict_putp(ctx, dict, path, val);
	}
	fz_always(ctx)
	{
		if (val)
			pdf_drop_obj(ctx, val);
	}
	fz_catch(ctx)
	{
		return 0;
	}

	return 1;
}

int safe_pdf_update_object(fz_context *ctx, pdf_document *pdf, int xref, pdf_obj *obj)
{
	fz_try(ctx)
	{
		pdf_update_object(ctx, pdf, xref, obj);
	}
	fz_catch(ctx)
	{
		return 0;
	}
	return 1;
}

fz_buffer *safe_pdf_load_stream_number(fz_context *ctx, pdf_document *pdf, int xref)
{
	fz_buffer *buf;
	fz_try(ctx)
	{
		buf = pdf_load_stream_number(ctx, pdf, xref);
	}
	fz_catch(ctx)
	{
		return NULL;
	}
	return buf;
}

fz_buffer *safe_pdf_load_raw_stream_number(fz_context *ctx, pdf_document *pdf, int xref)
{
	fz_buffer *buf;
	fz_try(ctx)
	{
		buf = pdf_load_raw_stream_number(ctx, pdf, xref);
	}
	fz_catch(ctx)
	{
		return NULL;
	}
	return buf;
}

int safe_pdf_dict_del_path(fz_context *ctx, pdf_obj *dict, const char *path)
{
	fz_try(ctx)
	{
		const char *slash = strrchr(path, '/');
		if (!slash)
		{
			pdf_dict_dels(ctx, dict, path);
		}
		else
		{
			size_t parent_len = (size_t)(slash - path);
			const char *leaf = slash + 1;
			if (parent_len == 0 || !*leaf)
				fz_throw(ctx, FZ_ERROR_ARGUMENT, "bad path");
			char *parent = (char *)malloc(parent_len + 1);
			if (!parent)
				fz_throw(ctx, FZ_ERROR_SYSTEM, "cannot allocate parent path");
			memcpy(parent, path, parent_len);
			parent[parent_len] = 0;
			pdf_obj *parent_obj = pdf_dict_getp(ctx, dict, parent);
			if (parent_obj && pdf_is_dict(ctx, parent_obj))
				pdf_dict_dels(ctx, parent_obj, leaf);
			free(parent);
		}
	}
	fz_catch(ctx)
	{
		return 0;
	}
	return 1;
}
*/
import "C"
