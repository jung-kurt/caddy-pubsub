BEGIN { 
	show = 0
	print "/*" 
}

/^\-/ { trim = 1 }

/^Package/ { show = 1 }

!NF { trim = 0 }

trim { sub("^ +", "", $0) }

show { print $0 }

# NR == 1, /^Package/ { print $0 }

END { print "*/\npackage pubsub" }

