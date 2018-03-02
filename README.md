# lrt
lavaorg runtime (lrt) common code appropriate for multiple projects

lrt provides a common set of utility code useful for other projects. the lrt will not perform any 'go get' of other packages.

The lrt packages will be self contained and only make use of the standard go library. The goal is a client applications only
need to rely on this group of packages and there will be no inadvertant extra packages brought in.

Where packages are sourced from external resources the source code will be incorporated into lrt and assume a package name
under the lrt namespace (appropriate attribution will be maintained to the original source). External packages adopted may
be modifed in significant ways.  

# Licensing
The entire project license is specified in the file LICENSE in the top level of the repository. (Apache 2)

Individual sub-directores or files may override with their own specified license.
