#ifndef GOPHERPDF_FITZ_CGO_H
#define GOPHERPDF_FITZ_CGO_H

#include <mupdf/fitz.h>
#include <stdlib.h>
#include <string.h>

extern const char *fz_version;
#if defined(_WIN32)
typedef unsigned long long store;
#else
typedef unsigned long store;
#endif

typedef struct pdf_document pdf_document;
typedef struct pdf_obj pdf_obj;
pdf_document *pdf_document_from_fz_document(fz_context *ctx, fz_document *doc);
int pdf_has_unsaved_changes(fz_context *ctx, pdf_document *doc);
int pdf_can_be_saved_incrementally(fz_context *ctx, pdf_document *doc);
int pdf_doc_was_linearized(fz_context *ctx, pdf_document *doc);
int pdf_was_repaired(fz_context *ctx, pdf_document *doc);

pdf_obj *pdf_trailer(fz_context *ctx, pdf_document *doc);
pdf_obj *pdf_dict_gets(fz_context *ctx, pdf_obj *dict, const char *key);
pdf_obj *pdf_dict_getp(fz_context *ctx, pdf_obj *dict, const char *path);
pdf_obj *pdf_dict_get(fz_context *ctx, pdf_obj *dict, pdf_obj *key);
void pdf_dict_putp(fz_context *ctx, pdf_obj *dict, const char *path, pdf_obj *val);
void pdf_dict_dels(fz_context *ctx, pdf_obj *dict, const char *key);
void pdf_update_object(fz_context *ctx, pdf_document *doc, int num, pdf_obj *newobj);
pdf_obj *pdf_new_name(fz_context *ctx, const char *name);
pdf_obj *pdf_load_name_tree(fz_context *ctx, pdf_document *doc, pdf_obj *which);
int pdf_is_dict(fz_context *ctx, pdf_obj *obj);
int pdf_dict_len(fz_context *ctx, pdf_obj *obj);
pdf_obj *pdf_dict_get_key(fz_context *ctx, pdf_obj *dict, int i);
pdf_obj *pdf_dict_get_val(fz_context *ctx, pdf_obj *dict, int i);
void pdf_dict_put_val_null(fz_context *ctx, pdf_obj *dict, pdf_obj *key);
int pdf_is_name(fz_context *ctx, pdf_obj *obj);
int pdf_is_string(fz_context *ctx, pdf_obj *obj);
int pdf_is_bool(fz_context *ctx, pdf_obj *obj);
int pdf_is_real(fz_context *ctx, pdf_obj *obj);
const char *pdf_to_name(fz_context *ctx, pdf_obj *obj);
const char *pdf_to_text_string(fz_context *ctx, pdf_obj *obj);
int pdf_to_bool(fz_context *ctx, pdf_obj *obj);
int pdf_is_indirect(fz_context *ctx, pdf_obj *obj);
pdf_obj *pdf_resolve_indirect(fz_context *ctx, pdf_obj *ref);
int pdf_is_array(fz_context *ctx, pdf_obj *obj);
int pdf_array_len(fz_context *ctx, pdf_obj *arr);
pdf_obj *pdf_array_get(fz_context *ctx, pdf_obj *arr, int i);
int pdf_is_int(fz_context *ctx, pdf_obj *obj);
int pdf_is_number(fz_context *ctx, pdf_obj *obj);
int pdf_is_null(fz_context *ctx, pdf_obj *obj);
int pdf_to_int(fz_context *ctx, pdf_obj *obj);
int pdf_to_num(fz_context *ctx, pdf_obj *obj);
float pdf_to_real(fz_context *ctx, pdf_obj *obj);
int pdf_lookup_page_number(fz_context *ctx, pdf_document *doc, pdf_obj *page);
pdf_obj *pdf_lookup_page_obj(fz_context *ctx, pdf_document *doc, int needle);
int pdf_obj_num_is_stream(fz_context *ctx, pdf_document *doc, int num);
int pdf_xref_len(fz_context *ctx, pdf_document *doc);
pdf_obj *pdf_load_object(fz_context *ctx, pdf_document *doc, int num);
fz_buffer *pdf_load_stream_number(fz_context *ctx, pdf_document *doc, int num);
fz_buffer *pdf_load_raw_stream_number(fz_context *ctx, pdf_document *doc, int num);
void pdf_drop_obj(fz_context *ctx, pdf_obj *obj);
void pdf_print_obj(fz_context *ctx, fz_output *out, pdf_obj *obj, int tight, int ascii);
void pdf_set_page_labels(fz_context *ctx, pdf_document *doc, int index, int style, const char *prefix, int start);
void pdf_delete_page_labels(fz_context *ctx, pdf_document *doc, int index);

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

static inline fz_document *open_document(fz_context *ctx, const char *filename)
{
	fz_document *doc;

	fz_try(ctx)
	{
		doc = fz_open_document(ctx, filename);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return doc;
}

static inline fz_document *open_document_with_stream(fz_context *ctx, const char *magic, fz_stream *stream)
{
	fz_document *doc;

	fz_try(ctx)
	{
		doc = fz_open_document_with_stream(ctx, magic, stream);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return doc;
}

static inline fz_page *load_page(fz_context *ctx, fz_document *doc, int number)
{
	fz_page *page;

	fz_try(ctx)
	{
		page = fz_load_page(ctx, doc, number);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return page;
}

static inline int run_page_contents(fz_context *ctx, fz_page *page, fz_device *dev, fz_matrix transform, fz_cookie *cookie)
{
	fz_try(ctx)
	{
		fz_run_page_contents(ctx, page, dev, transform, cookie);
	}
	fz_catch(ctx)
	{
		return 0;
	}

	return 1;
}

static inline int run_page_annots(fz_context *ctx, fz_page *page, fz_device *dev, fz_matrix transform, fz_cookie *cookie)
{
	fz_try(ctx)
	{
		fz_run_page_annots(ctx, page, dev, transform, cookie);
	}
	fz_catch(ctx)
	{
		return 0;
	}

	return 1;
}

static inline fz_outline *load_outline(fz_context *ctx, fz_document *doc)
{
	fz_outline *outline;

	fz_try(ctx)
	{
		outline = fz_load_outline(ctx, doc);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return outline;
}

static inline char *copy_rectangle(fz_context *ctx, fz_stext_page *page, fz_rect area, int crlf)
{
	char *text;

	fz_try(ctx)
	{
		text = fz_copy_rectangle(ctx, page, area, crlf);
	}
	fz_catch(ctx)
	{
		return NULL;
	}

	return text;
}

static void append_name_entry(fz_context *ctx, pdf_document *pdf, fz_buffer *buf, const char *name, pdf_obj *val)
{
	pdf_obj *dest = val;
	if (pdf_is_indirect(ctx, dest))
		dest = pdf_resolve_indirect(ctx, dest);

	if (pdf_is_dict(ctx, dest))
		dest = pdf_dict_gets(ctx, dest, "D");

	if (!pdf_is_array(ctx, dest))
		return;

	int n = pdf_array_len(ctx, dest);
	if (n < 2)
		return;

	int page = -1;
	float x = 0;
	float y = 0;
	float z = 0;

	pdf_obj *page_obj = pdf_array_get(ctx, dest, 0);
	if (page_obj)
	{
		if (pdf_is_indirect(ctx, page_obj) || pdf_is_dict(ctx, page_obj))
		{
			page = pdf_lookup_page_number(ctx, pdf, page_obj);
		}
		else if (pdf_is_int(ctx, page_obj))
		{
			page = pdf_to_int(ctx, page_obj);
		}
	}

	pdf_obj *kind_obj = pdf_array_get(ctx, dest, 1);
	if (kind_obj && pdf_is_name(ctx, kind_obj))
	{
		const char *kind = pdf_to_name(ctx, kind_obj);
		if (kind && strcmp(kind, "XYZ") == 0)
		{
			if (n > 2)
			{
				pdf_obj *xo = pdf_array_get(ctx, dest, 2);
				if (xo && pdf_is_number(ctx, xo))
					x = pdf_to_real(ctx, xo);
			}
			if (n > 3)
			{
				pdf_obj *yo = pdf_array_get(ctx, dest, 3);
				if (yo && pdf_is_number(ctx, yo))
					y = pdf_to_real(ctx, yo);
			}
			if (n > 4)
			{
				pdf_obj *zo = pdf_array_get(ctx, dest, 4);
				if (zo && pdf_is_number(ctx, zo))
					z = pdf_to_real(ctx, zo);
			}
		}
	}

	fz_append_printf(ctx, buf, "%s\t%d\t%.10g\t%.10g\t%.10g\n", name, page, x, y, z);
}

static void append_name_tree(fz_context *ctx, pdf_document *pdf, fz_buffer *buf, pdf_obj *dict)
{
	int count = pdf_dict_len(ctx, dict);
	for (int i = 0; i < count; i++)
	{
		pdf_obj *key = pdf_dict_get_key(ctx, dict, i);
		pdf_obj *val = pdf_dict_get_val(ctx, dict, i);
		if (!key || !val || !pdf_is_name(ctx, key))
			continue;
		const char *name = pdf_to_name(ctx, key);
		if (!name || !*name)
			continue;
		append_name_entry(ctx, pdf, buf, name, val);
	}
}

static inline char *resolve_names_dump(fz_context *ctx, fz_document *doc)
{
	pdf_document *pdf = pdf_document_from_fz_document(ctx, doc);
	if (!pdf)
		return NULL;

	fz_buffer *buf = NULL;
	char *out = NULL;

	fz_try(ctx)
	{
		buf = fz_new_buffer(ctx, 2048);
		pdf_obj *catalog = pdf_dict_gets(ctx, pdf_trailer(ctx, pdf), "Root");
		pdf_obj *dests_name = pdf_new_name(ctx, "Dests");

		pdf_obj *old_dests = pdf_dict_get(ctx, catalog, dests_name);
		if (old_dests && pdf_is_dict(ctx, old_dests))
			append_name_tree(ctx, pdf, buf, old_dests);

		pdf_obj *tree = pdf_load_name_tree(ctx, pdf, dests_name);
		if (tree && pdf_is_dict(ctx, tree))
			append_name_tree(ctx, pdf, buf, tree);

		size_t n = fz_buffer_storage(ctx, buf, NULL);
		out = (char *)malloc(n + 1);
		if (!out)
		{
			fz_throw(ctx, FZ_ERROR_SYSTEM, "cannot allocate resolve names output");
		}
		memcpy(out, fz_string_from_buffer(ctx, buf), n);
		out[n] = 0;
	}
	fz_always(ctx)
	{
		if (buf)
			fz_drop_buffer(ctx, buf);
	}
	fz_catch(ctx)
	{
		if (out)
		{
			free(out);
			out = NULL;
		}
	}

	return out;
}

static inline pdf_obj *safe_pdf_load_object(fz_context *ctx, pdf_document *pdf, int num)
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

static inline int safe_pdf_print_obj(fz_context *ctx, fz_output *out, pdf_obj *obj)
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

static inline int safe_pdf_dict_put_val_null(fz_context *ctx, pdf_obj *dict, pdf_obj *key)
{
	fz_try(ctx)
	{
		pdf_dict_put_val_null(ctx, dict, key);
	}
	fz_catch(ctx)
	{
		return 0;
	}
	return 1;
}

static inline int safe_pdf_form_field_count(fz_context *ctx, pdf_document *pdf)
{
	int count = -1;
	fz_try(ctx)
	{
		pdf_obj *trailer = pdf_trailer(ctx, pdf);
		pdf_obj *root = pdf_dict_gets(ctx, trailer, "Root");
		pdf_obj *acro = pdf_dict_gets(ctx, root, "AcroForm");
		pdf_obj *fields = pdf_dict_gets(ctx, acro, "Fields");
		if (fields && pdf_is_array(ctx, fields))
			count = pdf_array_len(ctx, fields);
	}
	fz_catch(ctx)
	{
		return -1;
	}
	return count;
}

static inline int safe_pdf_set_page_labels(fz_context *ctx, pdf_document *pdf, int index, int style, const char *prefix, int start)
{
	fz_try(ctx)
	{
		pdf_set_page_labels(ctx, pdf, index, style, prefix, start);
	}
	fz_catch(ctx)
	{
		return 0;
	}
	return 1;
}

static inline int safe_pdf_delete_page_labels(fz_context *ctx, pdf_document *pdf, int index)
{
	fz_try(ctx)
	{
		pdf_delete_page_labels(ctx, pdf, index);
	}
	fz_catch(ctx)
	{
		return 0;
	}
	return 1;
}

static inline pdf_obj *safe_pdf_parse_stm_obj_from_string(fz_context *ctx, pdf_document *pdf, const char *text)
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

static inline int safe_pdf_dict_putp_parsed(fz_context *ctx, pdf_document *pdf, pdf_obj *dict, const char *path, const char *value)
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

static inline int safe_pdf_update_object(fz_context *ctx, pdf_document *pdf, int xref, pdf_obj *obj)
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

static inline fz_buffer *safe_pdf_load_stream_number(fz_context *ctx, pdf_document *pdf, int xref)
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

static inline fz_buffer *safe_pdf_load_raw_stream_number(fz_context *ctx, pdf_document *pdf, int xref)
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

static inline int safe_pdf_dict_del_path(fz_context *ctx, pdf_obj *dict, const char *path)
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

static void append_page_label_rules_from_node(fz_context *ctx, pdf_obj *node, fz_buffer *buf)
{
	if (!node)
		return;
	if (pdf_is_indirect(ctx, node))
		node = pdf_resolve_indirect(ctx, node);
	if (!node || !pdf_is_dict(ctx, node))
		return;

	pdf_obj *nums = pdf_dict_gets(ctx, node, "Nums");
	if (nums && pdf_is_array(ctx, nums))
	{
		int n = pdf_array_len(ctx, nums);
		for (int i = 0; i + 1 < n; i += 2)
		{
			pdf_obj *idx = pdf_array_get(ctx, nums, i);
			pdf_obj *rule = pdf_array_get(ctx, nums, i + 1);
			if (!idx || !rule || !pdf_is_int(ctx, idx))
				continue;
			if (pdf_is_indirect(ctx, rule))
				rule = pdf_resolve_indirect(ctx, rule);
			if (!rule || !pdf_is_dict(ctx, rule))
				continue;

			int startpage = pdf_to_int(ctx, idx);
			const char *prefix = "";
			const char *style = "";
			int first = 1;

			pdf_obj *pobj = pdf_dict_gets(ctx, rule, "P");
			if (pobj && pdf_is_string(ctx, pobj))
				prefix = pdf_to_text_string(ctx, pobj);

			pdf_obj *sobj = pdf_dict_gets(ctx, rule, "S");
			if (sobj && pdf_is_name(ctx, sobj))
				style = pdf_to_name(ctx, sobj);

			pdf_obj *stobj = pdf_dict_gets(ctx, rule, "St");
			if (stobj && pdf_is_number(ctx, stobj))
				first = pdf_to_int(ctx, stobj);

			fz_append_printf(ctx, buf, "%d\t%s\t%s\t%d\n", startpage, style ? style : "", prefix ? prefix : "", first);
		}
		return;
	}

	pdf_obj *kids = pdf_dict_gets(ctx, node, "Kids");
	if (kids && pdf_is_array(ctx, kids))
	{
		int n = pdf_array_len(ctx, kids);
		for (int i = 0; i < n; i++)
		{
			pdf_obj *kid = pdf_array_get(ctx, kids, i);
			append_page_label_rules_from_node(ctx, kid, buf);
		}
	}
}

static inline char *page_label_rules_dump(fz_context *ctx, fz_document *doc)
{
	pdf_document *pdf = pdf_document_from_fz_document(ctx, doc);
	if (!pdf)
		return NULL;

	fz_buffer *buf = NULL;
	char *out = NULL;

	fz_try(ctx)
	{
		buf = fz_new_buffer(ctx, 1024);
		pdf_obj *trailer = pdf_trailer(ctx, pdf);
		pdf_obj *root = pdf_dict_gets(ctx, trailer, "Root");
		pdf_obj *labels = pdf_dict_gets(ctx, root, "PageLabels");
		append_page_label_rules_from_node(ctx, labels, buf);

		size_t n = fz_buffer_storage(ctx, buf, NULL);
		out = (char *)malloc(n + 1);
		if (!out)
		{
			fz_throw(ctx, FZ_ERROR_SYSTEM, "cannot allocate page-label output");
		}
		memcpy(out, fz_string_from_buffer(ctx, buf), n);
		out[n] = 0;
	}
	fz_always(ctx)
	{
		if (buf)
			fz_drop_buffer(ctx, buf);
	}
	fz_catch(ctx)
	{
		if (out)
		{
			free(out);
			out = NULL;
		}
	}

	return out;
}

#endif
