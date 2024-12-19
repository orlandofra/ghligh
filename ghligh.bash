#!/bin/bash

SUBCOMMANDS="export import cat ls"
COMMAND="ghligh"

_ghligh_completion() {
	local cur prev subcommand opts
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	subcommand="${COMP_WORDS[1]}"

	case "$subcommand" in
		export|import)
			opts="--from -f --to -t"
			if [[ "$prev" == "--from" || "$prev" == "-f" || "$prev" == "--to" || "$prev" == "-t" ]]; then
				COMPREPLY=( $(compgen -f -- "$cur") )
			else
				COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
			fi
			;;
		cat|ls)
			COMPREPLY=( $(compgen -f -- "$cur") )
			;;
		*)
			COMPREPLY=( $(compgen -W "$SUBCOMMANDS" -- "$cur") )
		;;
	esac
}

complete -o filenames -F _ghligh_completion $COMMAND

