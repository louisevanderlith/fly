#!/bin/bash

#==============================================================================
# bin script header:
# All scripts in bin directory should start with this
# this header is automatically added during the make process
#==============================================================================

#-----------------------------------------------------------------
# nr of lines in this header file
# this value is set during the make process
# and is used as an offset when reporting line numbers
# in debug statements, because this header is prefixed
# to script files, making the actual line number in the source
# .script file not the same as that of the released .sh file
#-----------------------------------------------------------------
BIN_SCRIPT_HEADER_LEN=66

#echo "Called $0 $*" >&2

#-----------------------------------------------------------------
# see if executing in release environment
#-----------------------------------------------------------------
tried_locations=""

p=$(dirname $0)
#echo "p=${p}" >&2
if [[ "$(basename $p)" != "libexec" && "$(basename $p)" != "bin" ]]
then
  p=$(dirname $p)
  #echo "p=${p}" >&2
fi      

script_setup="${p}/../libexec/bash/script_setup.bash"
if [ -f "${script_setup}" ]
then
  source ${script_setup}
  tried_locations=""
  debug "Executing in release environment." 
else
  #----------------------------------------------------
  # see if executing in release environment from component script
  #----------------------------------------------------
  tried_locations="Not found: ${script_setup}"
  script_setup="$(dirname $0)/../bash/script_setup.bash"
  if [ -f "${script_setup}" ]
  then
    source ${script_setup}
    tried_locations=""
    debug "Executing in dev environment."
  else
    tried_locations="${tried_locations}\nNot found: ${script_setup}"
    script_setup=""
  fi
fi

if [ -z "${script_setup}" ]
then
  echo "ERROR: Missing script_setup.bash" >&2
  echo -e "${tried_locations}" >&2
  exit 1
fi

#==============================================================================
# end of script header added during make process
#==============================================================================
#---------------------------------------------------------------------------
# DEV SCRIPTS: During make, a header is added and the file renamed to .sh
# This file is not a complete bash script and can only run after being build
#---------------------------------------------------------------------------
opt_arch="amd64"
opt_goos="linux"
opt_install=1
opt_package=0

COMMAND_LINE_BEGIN
COMMAND_LINE_OPTION -a opt_arch    STR      "Target Architecture"
COMMAND_LINE_OPTION -o opt_goos    STR      "Target Operating System"
COMMAND_LINE_OPTION -i opt_install NO_VALUE "Install (implied)"
COMMAND_LINE_OPTION -p opt_package NO_VALUE "Create tgz package file"
COMMAND_LINE_END $*

debug "Option: opt_arch=${opt_arch}"
debug "Option: opt_goos=${opt_goos}"
[ ${opt_package} -ne 0 ] && opt_install=1

#-------------------------------------------
# project directory must be git controlled
#-------------------------------------------
proj=$(git remote -v | grep "\.git " | head -n1 | awk '{print $2}' | xargs basename | sed "s/\.git//")
[ $? -ne 0 ] && error "Current directory is not version controlled with GIT."
[ -z "${proj}" ] && error "Cannot determine project name from git info, looking for xxx.git"
debug "proj=${proj}"

version=$(git describe --tags --always --dirty)
debug "version=${version}"

build=`git rev-parse HEAD`
debug "build=${build}"

#-----------------------------------------------------------------------------
# Setup the -ldflags option for go build here, interpolate the variable values
#-----------------------------------------------------------------------------
ldflags="-X main.ProductName=${proj} -X main.ServiceVersion=${version} -X main.Build=${build} -X env.ServiceVersion=${version}"
debug "ldflags=${ldflags}"

#-------------------------------------
# default VSERVICES_DIR if not defined
#-------------------------------------
[ -z "${VSERVICES_DIR}" ] && VSERVICES_DIR=${vdir}
[ -z "${VSERVICES_DIR}" ] && VSERVICES_DIR="/opt/vservices"

if [ ${opt_install} -ne 0 ]
then
	install_dir=${VSERVICES_DIR}/releases/${proj}-${version}-${opt_arch}
	debug "install_dir=${install_dir}"

	rm -fR ${install_dir}
	mkdir -p ${install_dir}/bin      || error "Failed to create ${install_dir}/bin"
	mkdir -p ${install_dir}/libexec  || error "Failed to create ${install_dir}/libexec"
	mkdir -p ${install_dir}/conf     || error "Failed to create ${install_dir}/conf"
	mkdir -p ${install_dir}/share    || error "Failed to create ${install_dir}/share"
fi

#----------------------------------------------------------
# keep errors in temp file so that whole build can complete
# and not stop on first error
#----------------------------------------------------------
error_file=$(mktemp)
error_count=0

#----------------------------------------------------------
# Go build in each sub-directories with a Go "package main"
# Excluding: "*/vendor/*"
#----------------------------------------------------------
main_dirs=$(find ./ -name "*.go" | xargs grep -l "package main" | xargs -n1 dirname 2>/dev/null | sort -u | grep -v -e "/vendor/")
for dir in ${main_dirs}
do
	pushd ${dir} >/dev/null
	verbose "Building Go executable in ${dir}"

	#-------------------------------------------
	# we want to allow build from any sub-directory
	# of the project. when called from project home
	# with exec in sub-dirs, this will have dir=...
	# where package main resides, but when called
	# from that directory, dir will be "." and
	# in this case, we only build the current
	# dir
	#-------------------------------------------
	# name of go binary is name of directory
	#-------------------------------------------
	binary=${install_dir}/libexec/$(basename ${dir})
	if [ "${binary}" == "." ]
	then
		binary=$(basename $(pwd))
	fi
	debug "  binary=${binary}"

	# delete existing binary 
	rm -f ${binary}

	# compile
	GOOS=${opt_goos} GOARCH=${opt_arch} go build -o ${binary} -ldflags "${ldflags}" 2>>${error_file}
	if [ $? -ne 0 ]
	then
		debug "Failed on ${binary} (goos:${opt_goos},arch:${opt_arch}). Adding to error report..."
		echo "***** BUILD FAILED for ${binary} *****" >>${error_file}
		let error_count=error_count+1
	else
		debug "Build completed for ${binary}"
		if [ ! -x ${binary} ]
		then
			echo "***** Missing binary $(pwd)/${binary} after go build *****" >> ${error_file}
			let error_count=error_count+1
		else
			# install (optional)
			if [ ${opt_install} -ne 0 ]
			then
				#JS: No need to copy, be output to that directory already
				#cp ${binary} ${install_dir}/libexec/ || error "Failed to copy ${binary} to ${install_dir}/bin"

				#------------------------------------------------------------------------
				# also copy conf files
				# all images copy to same target dir, so make sure no filename duplicates
				#------------------------------------------------------------------------
				if [ -d ./conf ]
				then
					ls -1 ./conf\
					| while read name
					do
						[ -f ${install_dir}/conf/${name} ] && debug "ERROR: ${install_dir}/conf/${name} already exist."
						debug "Installing ${install_dir}/conf/${name} ok"
					done
					cp -r conf/* ${install_dir}/conf/
				fi
				verbose "${binary}: Installed ${install_dir}/${binary}"
			else
				verbose "${binary}: Built successfully"
			fi
		fi
	fi
	popd >/dev/null
done

#---------------------------------------------------------------
# use automake for non-Go code...
#---------------------------------------------------------------
configure_files=$(find ./ -name "configure.ac")
for configure_file in ${configure_files}
do
	pushd $(dirname ${configure_file})

	if [ ! -f ./configure ]
	then
			aclocal || error "Failed in aclocal"
			autoconf || error "Failed in autoconf"
			automake --add-missing || error "Failed in automake"
			debug "Generated configure script and make files"
	fi

	if [ ${opt_install} -eq 0 ]
	then
		verbose "$(pwd): autoconf..."
		./configure >> ${error_file} 2>&1
	else
		verbose "$(pwd): autoconf -> ${install_dir}..."
		./configure --prefix=${install_dir} >> ${error_file} 2>&1
	fi

	if [ $? -ne 0 ]
	then
		echo "***** $(pwd) Auto configure failed *****" >> ${error_file}
		let error_count=error_count+1
	fi
	
	if [ ${opt_install} -eq 0 ]
	then
		verbose "$(pwd): automake..."
		make >> ${error_file} >> ${error_file} 2>&1
	else
		verbose "$(pwd): automake -> ${install_dir}..."
		make install >> ${error_file} >> ${error_file} 2>&1
	fi

	if [ $? -ne 0 ]
	then
		echo "***** $(pwd) Auto make failed *****" >> ${error_file}
		let error_count=error_count+1
	fi

	popd
done

if [ ${opt_install} -ne 0 ]
then
	#-----------------------------------------------------
	# create link to this latest release in releases
	# so one can test locally with the latest code without
	# having to switch after each commit
	#-----------------------------------------------------
	rm -f ${VSERVICES_DIR}/releases/${proj}
	ln -s ${install_dir} ${VSERVICES_DIR}/releases/${proj}

	verbose ""
	verbose "New release: $(ls -l ${VSERVICES_DIR}/releases/${proj})"
fi

if [ ${error_count} -gt 0 ]
then
	echo "=====[ BUILD ERRORS ]===================================================="
	cat ${error_file}
	rm -f ${error_file}
	echo "=====[ BUILD FAILED ]===================================================="
	exit 1
fi

if [ ${opt_package} -ne 0 ]
then
	#-----------------------------------------------
	# create package file from the install directory
	#-----------------------------------------------
	package_file=${install_dir}.tar.gz
	pushd $(dirname ${install_dir}) >/dev/null
	tar -czvf ${package_file} $(basename ${install_dir})

	verbose ""
	verbose "Packaged into ${package_file}:"
	ls -l ${package_file}
	popd >/dev/null
fi

rm -f ${error_file}
verbose "Build successfully."
exit 0
