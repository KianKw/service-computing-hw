/*=================================================================

Program name:
	selpg (SELect PaGes)

Purpose:
	Sometimes one needs to extract only a specified range of
pages from an input text file. This program allows the user to do
that.

Author: Kian Kwok

===================================================================*/

/*================================= includes ======================*/

package main

import (
	"fmt"
	"io"
	"os"
	"bufio"
	"strconv"
	"os/exec"
	flag "github.com/spf13/pflag"
)

/*================================= types =========================*/

type selpg_args struct {
	start_page int
	end_page int
	in_filename string
	page_len int /* default value, can be overriden by "-l number" on command line */
	page_type int /* 'l' for lines-delimited, 'f' for form-feed-delimited */
					/* default is 'l' */
	print_dest string
}

/*================================= globals =======================*/

var progname string /* program name, for error messages */

/*================================= prototypes ====================*/

/*================================= main() ========================*/

func main() {
	sa := selpg_args{}

	/* save name by which program is invoked, for error messages */
	progname = os.Args[0]

	flag.IntVarP(&sa.start_page, "start_page", "s", 0, "Input the start page")
	flag.IntVarP(&sa.end_page, "end_page", "e",  0, "Input the end page")
	flag.StringVarP(&sa.in_filename, "in_filename", "", "", "Input the in_filename");
	flag.IntVarP(&sa.page_len, "page_len", "l", 10, "Input the lines per page")
	flag.IntVarP(&sa.page_type, "page_type", "f", 'l', "Input the page type, 'l' for line-delimited, 'f' for for,-feed-delimited, default is 'l'")
	flag.StringVarP(&sa.print_dest, "print_dest", "d", "", "Input the print dest file")
	flag.Parse()
	
	process_args(&sa)
	process_input(sa)
}

/*================================= process_args() ================*/

func process_args(psa * selpg_args) {
	ac := len(os.Args)
	var s1 string /* temp str */
	var i int

	/* check the command-line arguments for validity */
	if ac < 3 {	/* Not enough args, minimum command is "selpg -sstartpage -eend_page"  */
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}

	/* handle mandatory args first */
	/* handle 1st arg - start page */
	s1 = os.Args[1] /* !!! PBO */
	if s1 != "-s" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s start_page\n", progname)
		usage()
		os.Exit(2)
	}
	INT_MAX := 1 << 32
	s1 = os.Args[2]
	i, _ = strconv.Atoi(s1)
	if i < 1 || i > (INT_MAX - 1) {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %s\n", progname, psa.start_page)
		usage()
		os.Exit(3)
	}
	psa.start_page = i

	/* handle 2nd arg - end page */
	s1 = os.Args[3] /* !!! PBO */
	if s1 != "-e" {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -e end_page\n", progname)
		usage()
		os.Exit(4)
	}
	s1 = os.Args[4]
	i, _ = strconv.Atoi(s1)
	if i < 1 || i > (INT_MAX - 1) || i < psa.start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %s\n", progname, psa.end_page)
		usage()
		os.Exit(5)
	}

	/* now handle optional args */
	/* handle page_len */
	if psa.page_len < 1 || psa.page_len > (INT_MAX - 1) {
		fmt.Fprintf(os.Stderr, "%s: invalid page length %s\n", progname, psa.page_len)
		usage()
		os.Exit(6)
	}
	
	/* handle in_filename */ 
	if len(flag.Args()) == 1 { /* there is one more arg */
		_, err := os.Stat(flag.Args()[0])
		/* check if file exists */
		if err != nil && os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n",
					progname, flag.Args()[0]);
			os.Exit(7);
		}
		psa.in_filename = flag.Args()[0]
	}
}

/*================================= process_input() ===============*/

func process_input(sa selpg_args) {
	var fin *os.File /* input stream */
	var fout io.WriteCloser /* output stream */
	var page_ctr int /* page counter */
	var line_ctr int /* line counter */

	/* set the input source */
	if len(sa.in_filename) == 0 {
		fin = os.Stdin
	} else {
		var err error
		fin, err = os.Open(sa.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n",
				progname, sa.in_filename)
			os.Exit(8)
		}
		defer fin.Close()
	}

	/* use  bufio.NewReader() to set a big buffer for fin, for performance */
	bufFin := bufio.NewReader(fin)

	/* set the output destination */
	cmd := &exec.Cmd{}
	if len(sa.print_dest) == 0 {
		fout = os.Stdout
	} else {
		cmd = exec.Command("cat")	
		var err error
		cmd.Stdout, err = os.OpenFile(sa.print_dest, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open file %s\n",
				progname, sa.print_dest)
			os.Exit(9)
		}

		fout, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to file %s\n",
				progname, sa.print_dest)
			os.Exit(10)
		}
		
		cmd.Start()
		defer fout.Close()
	}

	/* begin one of two main loops to print result based on page type */
	if sa.page_type == 'l' {
		line_ctr = 0
		page_ctr = 1
		for {
			line,  err := bufFin.ReadString('\n')
			if err != nil {	/* error or EOF */
				break
			}
			line_ctr++
			if line_ctr > sa.page_len {
				page_ctr++
				line_ctr = 1
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				_, err := fout.Write([]byte(line))
				if err != nil {
					fmt.Println(err)
					os.Exit(11)
				}
		 	}
		}  
	} else {
		page_ctr = 1
		for {
			page, err := bufFin.ReadString('\f')
			if err != nil { /* error or EOF */
				break
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				_, err := fout.Write([]byte(page))
				if err != nil {
					os.Exit(12)
				}
			}
			page_ctr++
		}
	}

	/* end main loop */
	
	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr,
			"%s: start_page (%d) greater than total pages (%d)," +
			" no output written\n", progname, sa.start_page, page_ctr)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr,"%s: end_page (%d) greater than total pages (%d)," +
		" less output than expected\n", progname, sa.end_page, page_ctr)
	}
}

/*================================= usage() =======================*/

func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s -s start_page -e end_page [-f|-l lines_per_page] [-d dest] [in_filename]\n", progname);
	flag.PrintDefaults()
}

/*================================= EOF ===========================*/
