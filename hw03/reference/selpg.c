/*=================================================================

Program name:
	selpg (SELect PaGes)

Purpose:
	Sometimes one needs to extract only a specified range of
pages from an input text file. This program allows the user to do
that.

Author: Vasudev Ram

===================================================================*/

/*================================= includes ======================*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <limits.h>
#include <assert.h>
#include <errno.h>

/*================================= types =========================*/

struct selpg_args
{
	int start_page;
	int end_page;
	char in_filename[BUFSIZ];
	int page_len; /* default value, can be overriden by "-l number" on command line */
	int page_type; /* 'l' for lines-delimited, 'f' for form-feed-delimited */
					/* default is 'l' */
	char print_dest[BUFSIZ];
};
typedef struct selpg_args sp_args;

/* INBUFSIZ is size of array inbuf */
#define INBUFSIZ ( 16 * 1024 )

/*================================= globals =======================*/

char *progname; /* program name, for error messages */

/*================================= prototypes ====================*/

void usage(void);
int main(int ac, char **av);
void process_args(int ac, char **av, sp_args* psa);
void process_input(sp_args sa);

/*================================= main()=== =====================*/

int main(int ac, char **av)
{
	sp_args sa;

	/* save name by which program is invoked, for error messages */
	progname = av[0];

	sa.start_page = sa.end_page = -1;
	sa.in_filename[0] = '\0';
	sa.page_len = 72;
	sa.page_type = 'l';
	sa.print_dest[0] = '\0';

	process_args(ac, av, &sa);
	process_input(sa);

	return 0;
}

/*================================= process_args() ================*/

void process_args(int ac, char **av, sp_args* psa)
{
	char s1[BUFSIZ]; /* temp str */
	char s2[BUFSIZ]; /* temp str */
	int argno; /* arg # currently being processed */
	/* arg at index 0 is the command name itself (selpg),
	   first actual arg is at index 1,
	   last arg is at index (ac - 1) */
	int i;

	/* check the command-line arguments for validity */
	if (ac < 3) /* Not enough args, minimum command is "selpg -sstartpage -eend_page"  */
	{
		fprintf(stderr, "%s: not enough arguments\n", progname);
		usage();
		exit(1);
	}

	/* handle mandatory args first */

#ifdef DEBUG
	fprintf(stderr, "DEBUG: before handling 1st arg\n");
#endif

	/* handle 1st arg - start page */
	strcpy(s1, av[1]); /* !!! PBO */
	if (strncmp(s1, "-s", 2) != 0)
	{
		fprintf(stderr, "%s: 1st arg should be -sstart_page\n", progname);
		usage();
		exit(2);
	}
	i = atoi(&s1[2]);
	if ( i < 1 || i > (INT_MAX - 1) )
	{
		fprintf(stderr, "%s: invalid start page %s\n", progname, &s1[2]);
		usage();
		exit(3);
	}
	psa->start_page = i;

#ifdef DEBUG
	fprintf(stderr, "DEBUG: before handling 2nd arg\n");
#endif

	/* handle 2nd arg - end page */
	strcpy(s1, av[2]); /* !!! PBO */
	if (strncmp(s1, "-e", 2) != 0)
	{
		fprintf(stderr, "%s: 2nd arg should be -eend_page\n", progname);
		usage();
		exit(4);
	}
	i = atoi(&s1[2]);
	if ( (i < 1) || (i > (INT_MAX - 1)) || (i < psa->start_page) )
	{
		fprintf(stderr, "%s: invalid end page %s\n", progname, &s1[2]);
		usage();
		exit(5);
	}
	psa->end_page = i;

	/* now handle optional args */

#ifdef DEBUG
	fprintf(stderr, "DEBUG: before while loop for opt args\n");
#endif

	argno = 3;
	while (argno <= (ac - 1) && av[argno][0] == '-')
	/* while there more args and they start with a '-' */
	{
		strcpy(s1, av[argno]); /* !!! PBO */
		switch(s1[1])
		{
			case 'l':
				strcpy(s2, &s1[2]); /* !!! PBO */
				i = atoi(s2);
				if ( i < 1 || i > (INT_MAX - 1) )
				{
					fprintf(stderr, "%s: invalid page length %s\n", progname, s2);
					usage();
					exit(6);
				}
				psa->page_len = i;
				++argno;
				continue;
				break; /* break after continue not reqd, defensive prog. */

			case 'f':
				/* check if just "-f" or something more */
				if (strcmp(s1, "-f") != 0)
				{
					fprintf(stderr, "%s: option should be \"-f\"\n", progname);
					usage();
					exit(7);
				}
				psa->page_type = 'f';
				++argno;
				continue;
				break;

			case 'd':
				strcpy(s2, &s1[2]); /* !!! PBO */
				/* check if dest specified */
				if (strlen(s2) < 1)
				{
					fprintf(stderr,
					"%s: -d option requires a printer destination\n", progname);
					usage();
					exit(8);
				}
				strcpy(psa->print_dest, s2); /* !!! PBO */
				++argno;
				continue;
				break;

			default:
				fprintf(stderr, "%s: unknown option %s\n", progname, s1);
				usage();
				exit(9);
				break;

		} /* end switch */
	} /* end while */

#ifdef DEBUG
	fprintf(stderr, "DEBUG: before check for filename arg\n");
	fprintf(stderr, "DEBUG: argno = %d\n", argno);
#endif

	/*++argno;*/
	if (argno <= (ac - 1)) /* there is one more arg */
	{
		strcpy(psa->in_filename, av[argno]); /* !!! PBO */
		/* check if file exists */
		if (access(psa->in_filename, F_OK) != 0)
		{
			fprintf(stderr, "%s: input file \"%s\" does not exist\n",
			progname, psa->in_filename);
			exit(10);
		}
		/* check if file is readable */
		if (access(psa->in_filename, R_OK) != 0)
		{
			fprintf(stderr, "%s: input file \"%s\" exists but cannot be read\n",
			progname, psa->in_filename);
			exit(11);
		}
	}

	assert(psa->start_page > 0);
	assert(psa->end_page > 0 && psa->end_page >= psa->start_page);
	assert(psa->page_len > 1);
	assert(psa->page_type == 'l' || psa->page_type == 'f');

#ifdef DEBUG
	fprintf(stderr, "DEBUG: psa->start_page = %d\n", psa->start_page);
	fprintf(stderr, "DEBUG: psa->end_page = %d\n", psa->end_page);
	fprintf(stderr, "DEBUG: psa->page_len = %d\n", psa->page_len);
	fprintf(stderr, "DEBUG: psa->page_type = %c\n", psa->page_type);
	fprintf(stderr, "DEBUG: psa->print_dest = %s\n", psa->print_dest);
	fprintf(stderr, "DEBUG: psa->in_filename = %s\n", psa->in_filename);
#endif

}

/*================================= process_input() ===============*/

void process_input(sp_args sa)
{
	FILE *fin; /* input stream */
	FILE *fout; /* output stream */
	char s1[BUFSIZ]; /* temp string var */
	char *crc; /* for char ptr return code */
	int c; /* to read 1 char */
	char line[BUFSIZ];
	int line_ctr; /* line counter */
	int page_ctr; /* page counter */
	char inbuf[INBUFSIZ]; /* for better performance on input stream */

	/* set the input source */
	if (sa.in_filename[0] == '\0')
	{
		fin = stdin;
	}
	else
	{
		fin = fopen(sa.in_filename, "r");
		if (fin == NULL)
		{
			fprintf(stderr, "%s: could not open input file \"%s\"\n",
			progname, sa.in_filename);
			perror("");
			exit(12);
		}
	}

	/* use setvbuf() to set a big buffer for fin, for performance */
	setvbuf(fin, inbuf, _IOFBF, INBUFSIZ);

	/* set the output destination */
	if (sa.print_dest[0] == '\0')
	{
		fout = stdout;
	}
	else
	{
		fflush(stdout);
		sprintf(s1, "lp -d%s", sa.print_dest);
		fout = popen(s1, "w");
		if (fout == NULL)
		{
			fprintf(stderr, "%s: could not open pipe to \"%s\"\n",
			progname, s1);
			perror("");
			exit(13);
		}
	}

	/* begin one of two main loops based on page type */
	if (sa.page_type == 'l')
	{
		line_ctr = 0;
		page_ctr = 1;

		while (1)
		{
			crc = fgets(line, BUFSIZ, fin);
			if (crc == NULL) /* error or EOF */
				break;
			line_ctr++;
			if (line_ctr > sa.page_len)
			{
				page_ctr++;
				line_ctr = 1;
			}
			if ( (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) )
			{
				fprintf(fout, "%s", line);
			}
		}
	}
	else
	{
		page_ctr = 1;
		while (1)
		{
			c = getc(fin);
			if (c == EOF) /* error or EOF */
				break;
			if (c == '\f') /* form feed */
				page_ctr++;
			if ( (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) )
			{
				putc(c, fout);
			}
		}
	}

	/* end main loop */

	if (page_ctr < sa.start_page)
	{
		fprintf(stderr,
		"%s: start_page (%d) greater than total pages (%d),"
		" no output written\n", progname, sa.start_page, page_ctr);
	}
	else if (page_ctr < sa.end_page)
	{
		fprintf(stderr,"%s: end_page (%d) greater than total pages (%d),"
		" less output than expected\n", progname, sa.end_page, page_ctr);
	}
	if (ferror(fin)) /* fgets()/getc() encountered an error on stream fin */
	{
		strcpy(s1, strerror(errno)); /* !!! PBO */
		fprintf(stderr, "%s: system error [%s] occurred on input stream fin\n",
		progname, s1);
		fclose(fin);
		exit(14);
	}
	else /* it was EOF, not error */
	{
		fclose(fin);
		fflush(fout);
		if (sa.print_dest[0] != '\0')
		{
			pclose(fout);
		}
		fprintf(stderr, "%s: done\n", progname);
	}
}

/*================================= usage() =======================*/

void usage(void)
{
	fprintf(stderr,
	"\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ]"
	" [ -ddest ] [ in_filename ]\n", progname);
}

/*================================= EOF ===========================*/
