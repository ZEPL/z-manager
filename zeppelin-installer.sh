#!/bin/bash
#
# Install Apache Zeppelin with Spark\Hadoop version

{ # Prevent execution if this script was only partially downloaded

readonly ZEPPELIN_LATEST="0.5.0-incubating-SNAPSHOT"
readonly SPARK_LATEST='1.3.1'
readonly HADOOP_LATEST='2.4.0'

readonly product_zeppelin='Apache Zeppelin (incubating)'
readonly product_manager='Z-Manager'
readonly please_enter='Please enter'



readonly gdrive="http://googledrive.com/host/0B0q6JdSf3mukfkF4Y1dkTjdNdk9wVmNRalY5dWlOTUV1QlpITUFXZWJSdjRRNkNXcFRUQ1k"


err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}
log() { #TODO(alex): add 'debug' switch
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@"
}

E_BAD_CURL=101
E_INSTALL_EXISTS=102
E_BAD_MOVE=103
E_BAD_ARCHIVE=104
E_BAD_INPUT=105
E_BAD_RM=106
E_BAD_LOGIN=107
E_BAD_READ=108
E_BAD_UI_PARAMS=109
E_NO_CURL_WGET=110

#Interface between UI and Installer: 5 basic vars + 5 for resource manager
zeppelin_ver="${ZEPPELIN_LATEST}"
spark_ver="${SPARK_LATEST}"
hadoop_ver="${HADOOP_LATEST}"
spark_cluster_master_url=""


resource_manager="spark standalone" #or yarn, mesos
yarn_hadoop_home_conf_dir=""
yarn_spark_home="/usr/lib/spark"
mesos_native_java_lib=""
mesos_spark_executor_uri=""
zeppelin_port="8080"

##############################################################
# CLI part of the Manager
# Gets the configuration variables from the user
##############################################################

declare -a spark_versions=('1.3.1' '1.3.0' '1.4.0')
declare -a hadoopv_spark_1_3_0=('2.0.0-cdh4.7.1')
declare -a hadoopv_spark_1_3_1=('1.0.4' '2.4.0' '2.0.0-cdh4.7.1' '2.5.0-cdh5.3.0')
declare -a hadoopv_spark_1_4_0=('2.6.0')
declare -r hadoop_default='2.4.0'
declare -r spark_default='1.3.1'
declare -r spark_latestt='1.3.1'
declare -r spark_experimental='1.4.0'
declare -r persist_filename='.zep'
declare -r re_yn='^[Yy]$|^[yY][eE][sS]$|^[Nn]$|^[nN][oO]$'
declare -r re_y='^[Yy]$|^[yY][eE][sS]$'
declare -r re_n='^[Nn]$|^[nN][oO]$'
declare -r re_num='^[0-9]+$'

declare -i num_ui_params=9



# interface variables with installer
declare -x HADOOP_VERSION="${hadoop_default}"
declare -x SPARK_VERSION="${spark_default}"
declare -x SPARK_MASTER_URL=""

declare -x RESOURCE_MANAGER="Spark standalone"

# YARN interface with installer
declare -x YARN_HADOOP_HOME_CONF_DIR=""
declare -x YARN_SPARK_HOME=""

# Mesos interface with installer
declare -x MESOS_NATIVE_JAVA_LIB=""
declare -x MESOS_SPARK_EXECUTOR_URI=""

declare -x ZEPPELIN_PORT="8080"



# ask user for either default or advanced installation
get_installation_type() {
  local installation_type
  echo
  echo 'Select the type of Zeppelin installation:'
  echo '1. Default (Hadoop 2.4.0; Spark 1.3.1; local mode)'
  echo '   Good for quickstart:'
  echo '       does not require local Spark installation'
  echo '       does not require external Spark cluster'
  echo '2. Advanced (pick versions and adjust configuration to your cluster)'
  read installation_type </dev/tty

  #TODO(khalid): change check of values (1,2) with regular expression

  while [[ "${installation_type}" != "1" && "${installation_type}" != "2" ]]; do
    echo "${please_enter} either 1 or 2"
    read installation_type </dev/tty
  done
  return "${installation_type}"

}

# request Hadoop version from the user
get_hadoop_version() {
  echo
  echo 'Please select one of the supported Hadoop releases'
  local sparkv_index=$1
  local -i len=0
  #TODO(khalid): make one -case statement in the beginning and pass required array
  case "${sparkv_index}" in
    1)
      for i in $(seq 1 ${#hadoopv_spark_1_3_1[@]}); do
        echo "${i}. Hadoop ${hadoopv_spark_1_3_1[$i-1]}"
      done
      len="${#hadoopv_spark_1_3_1[@]}"
      ;;
    2)
      for i in $(seq 1 ${#hadoopv_spark_1_3_0[@]}); do
        echo "${i}. Hadoop ${hadoopv_spark_1_3_0[$i-1]}"
      done
      len="${#hadoopv_spark_1_3_0[@]}"
      ;;
    3)
      for i in $(seq 1 ${#hadoopv_spark_1_4_0[@]}); do
        echo "${i}. Hadoop ${hadoopv_spark_1_4_0[$i-1]}"
      done
      len="${#hadoopv_spark_1_4_0[@]}"
      ;;
    *)
      err 'Invalid version of Spark'
      exit "${E_BAD_INPUT}"
      ;;
  esac

  local hadoop_selection
  read hadoop_selection </dev/tty

  while ! [[ "${hadoop_selection}" =~ $re_num &&
             "${hadoop_selection}" -ge 1  &&
             "${hadoop_selection}" -le ${len} ]]; do
    echo "${please_enter} the number between 1 and ${len}"
    read hadoop_selection </dev/tty
  done

  local -i array_index=${hadoop_selection}-1
  return "${array_index}"
}

# request Spark version from the user
get_spark_version() {
  echo
  echo 'Please select one of the following Spark releases'
  for i in $(seq 1 ${#spark_versions[@]}); do
    if [[ "${spark_versions[$i-1]}" = "${spark_latestt}" ]]; then
      echo "${i}. Spark ${spark_versions[$i-1]} (latest)"
    elif [[ "${spark_versions[$i-1]}" = "${spark_experimental}" ]]; then
      echo "${i}. Spark ${spark_versions[$i-1]} (experimental)"
    else
      echo "${i}. Spark ${spark_versions[$i-1]}"
    fi
  done

  read spark_selection </dev/tty

  while ! [[ "${spark_selection}" =~ $re_num &&
             "${spark_selection}" -ge 1  &&
             "${spark_selection}" -le ${#spark_versions[@]} ]]; do
    echo "${please_enter} the number between 1 and ${#spark_versions[@]}"
    read spark_selection </dev/tty
  done

  return "${spark_selection}"
}






#print non-empty val with given label
print_if_not_empty() {
  local key="$1"
  local val="$2"
  if [[ -n "${val}" ]]; then
    echo "${key}: ${val}"
  fi
}

# requesting user whether to proceed to installation
install_confirmation() {
  echo 
  echo 'Configuration is finished now,'
  echo 'please review the options you have chosen before installation begins:'
  echo

  print_if_not_empty "Spark version" "$SPARK_VERSION"
  print_if_not_empty "Hadoop version" "$HADOOP_VERSION"

  if [ -z "$SPARK_MASTER_URL" ]; then
    echo "Spark master URL: None"
  else
    print_if_not_empty  "Spark master URL" "$SPARK_MASTER_URL"
  fi


  print_if_not_empty "Resource manager" "$RESOURCE_MANAGER"
  print_if_not_empty "Spark home" "$YARN_SPARK_HOME"
  print_if_not_empty "Hadoop home" "$YARN_HADOOP_HOME_CONF_DIR"
  print_if_not_empty "Mesos Java library" "$MESOS_NATIVE_JAVA_LIB"
  print_if_not_empty "Mesos Spark executor URI" "$MESOS_SPARK_EXECUTOR_URI"
  print_if_not_empty "Zeppelin port" "${ZEPPELIN_PORT}"

  echo
  echo 'y(es)/n(o) ?'
  read yn </dev/tty 

  while ! [[ "${yn}" =~ $re_yn ]]; do
    echo "${please_enter} y(es)/n(o)"
    read yn </dev/tty
  done
  echo
  case "${yn}" in
    [Yy] | [yY][eE][sS])
      echo "Initiating installation..." ;;
    [Nn] | [nN][oO])
      echo 'exiting...' && exit ;;
    *)
      err 'Invalid input'
      exit "${E_BAD_INPUT}";;
  esac

}

# persist config parameters into the ${persist_filename} in the current folder.
# overwrite previous file in case of successful installation.
persist () {

  if [[ ! "${#}" = "${num_ui_params}" ]]; then
    err "Number of installation parameters is different from ${num_ui_params}"
    exit "${E_BAD_UI_PARAMS}"
  fi

  i=0
  for arg in "$@" ; do
    if [[ "${i}" = 0 ]]; then
      echo "${arg}" > "${persist_filename}"
    else
      echo "${arg}" >> "${persist_filename}"
    fi
    i=$(($i + 1))
  done

}

# show user preferences on previous installation; ask for reuse
show_history () {
  echo
  echo 'Previous installation history exists in this folder'
  echo

  local -a user_params
  i=0
  while read line; do
    user_params[i]="${line}"
    i=$(($i + 1))
  done < "${persist_filename}"

  if [[ ! "${i}" = "${num_ui_params}" ]]; then
    err "Number of installation parameters is different from ${num_ui_params}"
    exit "${E_BAD_UI_PARAMS}"
  fi

  print_if_not_empty "Spark version" "${user_params[0]}"
  print_if_not_empty "Hadoop version" "${user_params[1]}"
  if [ -z "${user_params[2]}" ]; then
    echo "Spark master URL: None"
  else
    print_if_not_empty  "Spark master URL" "${user_params[2]}"
  fi
  print_if_not_empty "Resource manager" "${user_params[3]}"
  print_if_not_empty "Spark home" "${user_params[4]}"
  print_if_not_empty "Hadoop home" "${user_params[5]}"
  print_if_not_empty "Mesos Java library" "${user_params[6]}"
  print_if_not_empty "Mesos Spark executor URI" "${user_params[7]}"
  print_if_not_empty "Zeppelin port" "${user_params[8]}"


  echo
  echo "Would you like to reuse last intallation setting? y(es)/n(o)"
  read yn </dev/tty

  while ! [[ "${yn}" =~ $re_yn ]]; do
    echo "${please_enter} y(es)/n(o)"
    read yn </dev/tty
  done

  case "${yn}" in
    [Yy] | [yY][eE][sS])
      SPARK_VERSION="${user_params[0]}"
      HADOOP_VERSION="${user_params[1]}"
      SPARK_MASTER_URL="${user_params[2]}"
      RESOURCE_MANAGER="${user_params[3]}"
      YARN_SPARK_HOME="${user_params[4]}"
      YARN_HADOOP_HOME_CONF_DIR="${user_params[5]}"
      MESOS_NATIVE_JAVA_LIB="${user_params[6]}"
      MESOS_SPARK_EXECUTOR_URI="${user_params[7]}"
      ZEPPELIN_PORT="${user_params[8]}"

      return 0
      ;;
    [Nn] | [nN][oO]) return 1 ;;
    *) return 2 ;;
  esac

}

# choose type of resource manager
get_cluster_manager() {
  echo
  echo "Select cluster resource manager:"
  echo "1. Spark Standalone (default)"
  echo "2. YARN"
  echo "3. Mesos"
  declare -i local manager=1

  read option </dev/tty

  while ! [[ "${option}" =~ $re_num &&
             "${option}" -ge 1  &&
             "${option}" -le 3 ]]; do
    echo "${please_enter} the number between 1 and 3"
    read option </dev/tty
  done

  case "${option}" in
    1)
      RESOURCE_MANAGER="Spark standalone"
      return 1
      ;;
    2)
      RESOURCE_MANAGER="Yarn"
      return 2
      ;;
    3)
      RESOURCE_MANAGER="Mesos"
      return 3
      ;;
    *) return 4 ;;
  esac

}


# request spark master url from the user
get_master_url() {
  local resource_manager="$1"

  if [[ "${resource_manager}" == "1" ]] || [[ "${resource_manager}" == "3" ]]; then
    echo
    echo "${please_enter} Spark cluster master URL"
    echo 'e.g.  spark://<your-spark-master-ip>:7077 or mesos://<your-mesos-master-ip>:5050'
    read SPARK_MASTER_URL </dev/tty
  elif [[ "${resource_manager}" == "2" ]]; then
    SPARK_MASTER_URL="yarn-client"
  fi
}

# get extra cluster parameters from the user
get_extra_cluster_params() {
  local resource_manager="$1"

  case "${resource_manager}" in
    1)
      # Nothing to do for Spark standalone resource manager setting
      : ;;
    2)
      # YARN resource manager setting
      echo
      echo "${please_enter} path to the \$HADOOP_CONF_DIR, where Hadoop configuration files are"      echo "e.g. yarn-site.xml, usually \$HADOOP_HOME/etc/hadoop:"
      read hadoop_path </dev/tty
      YARN_HADOOP_HOME_CONF_DIR="${hadoop_path}"
      ;;
    3)
      # Mesos resource manager setting
      echo
      echo "${please_enter} path to the Mesos native library (e.g. \${MESOS_HOME}/lib/libmesos.so (or .dylib))"
      local mesos_lib_path=""
      read mesos_lib_path </dev/tty
      MESOS_NATIVE_JAVA_LIB="${mesos_lib_path}"
      local spark_bin="spark-${spark_ver}.tar.gz"
      echo
      echo "${please_enter} path to the Spark binary package:"
      echo 'If in local filesystem - should be same on all workers, otherwise remotely accessible to all workers'
      echo "i.e http://.../${spark_bin} , s3n://.../${spark_bin} or hdfs://.../${spark_bin}"
      local mesos_exec_uri=""
      read mesos_exec_uri </dev/tty
      MESOS_SPARK_EXECUTOR_URI="${mesos_exec_uri}"
      ;;
    *) echo "Wrong resource manager option" ;;

  esac

}

# ask for spark installation directory
get_spark_home() {
  local resource_manager="$1"
  local spark_home=""
  echo
  echo "${please_enter} path to the Spark installation (e.g. /usr/lib/spark):"
  read spark_home </dev/tty
  YARN_SPARK_HOME="${spark_home}"
}

# asks user for port to run Zeppelin (default 8080)
get_zeppelin_port() {

  if [[ "${BASH_VERSINFO[0]}" < 4  ]] ; then
    read -p "${please_enter} Zeppelin port number: (default: 8080): " zeppelin_port </dev/tty
    zeppelin_port="${zeppelin_port:-8080}"
  else
    read -p "${please_enter} Zeppelin port number: " -e -i 8080 zeppelin_port </dev/tty
  fi
  echo "${zeppelin_port}"
}

# ask new installation parameters from the user
get_new_settings() {

  get_installation_type
  local install_type="$?"
  local resource_manager_id="1"

  if [[ "${install_type}" = "2" ]]; then
    get_spark_version
    #actually not index, but index + 1
    local spark_index="$?"
    SPARK_VERSION="${spark_versions[${spark_index}-1]}"
    get_hadoop_version "${spark_index}"
    #this one is index though
    local hadoop_index="$?"
    case "${spark_index}" in
      1) HADOOP_VERSION="${hadoopv_spark_1_3_1[${hadoop_index}]}" ;;
      2) HADOOP_VERSION="${hadoopv_spark_1_3_0[${hadoop_index}]}" ;;
      3) HADOOP_VERSION="${hadoopv_spark_1_4_0[${hadoop_index}]}" ;;
      *)
        err 'Wrong Spark selection'
        exit "${E_BAD_INPUT}"
        ;;
    esac

    echo
    echo 'Do you want to configure external Spark cluster? y(es)/n(o)'
    read yn </dev/tty
    while ! [[ "${yn}" =~ $re_yn ]]; do
      echo "${please_enter} y(es)/n(o)"
      read yn </dev/tty
    done

    if [[ "${yn}" =~ $re_y ]]; then
      get_cluster_manager
      resource_manager_id="$?"
      get_master_url "${resource_manager_id}"
      get_spark_home
      get_extra_cluster_params "${resource_manager_id}"
    fi

    echo
    ZEPPELIN_PORT="$(get_zeppelin_port)"

  elif [[ "${install_type}" = "1" ]]; then
    # default settings are chosen, nothing to do
    :
  else
    err 'Wrong installation type'
    exit "${E_BAD_INPUT}"
  fi


}



##############################################################
# Get user parameters and save them into four global variables
# Globals:
#   SPARK_VERSION       - version of Spark
#   HADOOP_VERSION      - version of Hadoop (depending on Spark)
#   SPARK_MASTER_URL   - URL of external Spark cluster (if available)

# Arguments:
#   None
# Returns:
#   None
##############################################################
start_ui () {
  echo
  echo "Welcome to ${product_manager}!"
  local use_last_setting=1

  if [[ -f "${persist_filename}" ]]; then
    show_history
    use_last_setting="$?"
  fi

  if [[ "${use_last_setting}" = 0 ]]; then
    # yes, reusing last setting. parameters are set, nothing to do
    :
  elif [[ "${use_last_setting}" = 1 ]]; then
    get_new_settings
  else
    err "Incorrect last installation parameters"
    exit "${E_BAD_UI_PARAMS}"
  fi

  install_confirmation
  echo
}



##############################################################
# Installer
##############################################################


##############################################################
# Helper function to download given filename from public GDrive
# Does not retry for now.
# Globals:
#   gdrive     - url to public folder
# Arguments:
#   filename   - name of the file to download
# Returns:
#   None
##############################################################
util::download_if_not_exits() {
  local filename="$1"
  local src_url="${gdrive}/${filename}"

  if [[ -e "${filename}" ]]; then
    log "Skip downloading ${filename} as it exisits"
    return 0  
  fi

  curl -kL --progress-bar -O "${src_url}"
  if [[ "$?" -ne 0 ]]; then #TODO(alex): would be nice to retry
    err "Unable to download ${build_filename} from ${src_url}" >&2
    exit "${E_BAD_CURL}"
  fi
}

util::count() {
  if [[ "${COUNT}" = true ]] ; then
    local tid="UA-38575365-10"
    local cid
    cid="$(echo $$)"

    curl -ks -d "v=1&tid=${tid}&cid=${cid}&t=event&ec=Zeppelin&ea=Download&el=Manager&ev=1" 'http://www.google-analytics.com/collect' > /dev/null
  else
    log "Opt-in installation count"
  fi
}

#mv src to dst with error reporting
util::mv_or_exit() {
  local src="$1"
  local dst="$2"
  if ! mv "${src}" "${dst}" ; then
    err "Unable to move ${src} to ${dst}"
    exit "${E_BAD_MOVE}"
  fi
}

##############################################################
# Downloads vanilla build of Zeppelin with proper spark\hadoop
# Globals:
#   None
# Arguments:
#   build_filename - full name of the Zeppelin distributive
# Returns:
#   None
##############################################################
download_zeppelin() {
  local build_filename="$1"
  log "Downloading ${product_zeppelin} from ${build_filename}..."

  #TODO(alex): ~200mb, so re-try would make sense
  util::download_if_not_exits "${build_filename}"
  util::count
  log "Done"
}

##############################################################
# Unpacks given .tar.gz distributive if it does not exit
# Globals:
#   None
# Arguments:
#   distr         - Zeppelin distributive filename
#   zeppelin_path - zeppelin folder name
# Returns:
#   None
##############################################################
unpack_distr() {
  local distr="$1"
  local zeppelin_path="$2"
  log "Unpack ${distr} to ${zeppelin_path} ..."

  if [[ -d "${zeppelin_path}" ]] ; then
    log "${zeppelin_path} already exists, exiting"
    exit "${E_INSTALL_EXISTS}"
  fi

  if ! tar xzf "${distr}" ; then
    err "Unable to extract ${distr}"
    exit "${E_BAD_ARCHIVE}"
  fi

  log "Done"
}



##############################################################
# Configures interpreter
# Globals:
#   None
# Arguments:
#   zepepiln_path    - path to local Z install
#   spark_master_url - URL to the external Spark Master
#   spark_home       - (YARN) dir with Spark installation
#   spark_executor_uri (Mesos) path to Spark binary
#   hadoop_conf_dir  - path to yarn-site.xml dir
#   mesos_native_lib - path to libmesos.so

# Returns:
#   None
##############################################################
configure_interpreter() {
  local zepepiln_path="$1"
  local master_url="$2"
  local spark_home="$3"
  local spark_executor_uri="$4"
  local hadoop_conf_dir="$5"
  local mesos_native_lib="$6"


  # interpreter.json
  if [[ -n "${master_url}" ]] ; then
    log "Downloading Spark interpreter configuration file ..."
    util::download_if_not_exits 'interpreter.json'
    log "Done"
    log "Adding Spark interpreter with ${master_url} to the interpreter.json"

    template_interp_json 'interpreter.json' \
                         "zeppelin-$(whoami)" \
                         "${master_url}" \
                         "${spark_home}" \
                         "${spark_executor_uri}"
    util::mv_or_exit 'interpreter.json' "${zepepiln_path}/conf"
    log "Done"
  fi

  # zeppelin-env.sh
  local conf_file='zeppelin-env.sh'
  log "Downloading Zeppelin configuration file ..."
  util::download_if_not_exits "${conf_file}"
  log "Done"

  local python_path
  python_path=$(deduce_pyspark_path "${spark_home}")
  log "Pyspark found at ${python_path}"
  template_zeppelin_env 'zeppelin-env.sh' \
                        "${hadoop_conf_dir}" \
                        "${mesos_native_lib}" \
                        "${python_path}" \
                        "${zeppelin_port}"


  util::mv_or_exit "${conf_file}" "${zepepiln_path}/conf"
  log "Done"
}

# replaces vars in interpreter.json template
template_interp_json() {
  local filename="$1"
  local app_name="$2"
  local master_url="$3"
  local spark_home="$4"
  local spark_executor_uri="$5"

  log "  Reading the ${filename}"
  local file1
  file1=$(<"${filename}")
  file1=${file1/<%spark_master_url%>/"${master_url}"}
  file1=${file1/<%spark_app_name%>/"${app_name}"}
  file1=${file1/<%spark_home%>/"${spark_home}"}
  file1=${file1/<%spark_executor_uri%>/"${spark_executor_uri}"}

  echo "${file1}" > "${filename}"
}

# guess PYTHONPATH needed for pyspark
deduce_pyspark_path() {
  local spark="$1"

  local zipfile=""
  #TODO(alex): replace with
  #zipfile="$(find "${yarn_spark_home}/python/lib" "-name" "py4j-*")"
  zipfile="${spark}/python/lib/py4j-0.8.2.1-src.zip"
  echo "${spark}/python:${zipfile}"
}

# replaces vars in zeppelin-env.sh template
template_zeppelin_env() {
  local filename="$1"
  local hadoop_home_conf_dir="$2"
  local mesos_native_java_lib="$3"
  local python_path="$4"
  
  log "  Reading the ${filename}"
  local file1
  file1=$(<"${filename}")
  file1=${file1/<%hadoop_home_conf_dir%>/"${hadoop_home_conf_dir}"}
  file1=${file1/<%mesos_native_java_lib%>/"${mesos_native_java_lib}"}
  file1=${file1/<%python_path%>/"${python_path}"}

  echo "${file1}" > "${filename}"
  echo "export ZEPPELIN_PORT=${zeppelin_port}" >> "${filename}"
}



# Prints out instructions on how to run Zeppelin
print_instructions_to_run() {
  echo
  echo "To run ${product_zeppelin} now do:"
  echo " ./bin/zeppelin-daemon.sh start"
  echo " and visit http://localhost:${zeppelin_port}"
}


COUNT='true'
while [[ $# > 0 ]] ; do
  key="$1"
  case $key in
    -n|--no-count)
    COUNT='false'
    shift
    ;;
    -s|--spark-version)
    spark_ver="$2"
    shift
    ;;
    -h|--hadoop-version)
    hadoop_ver="$2"
    shift
    ;;
    --default)
    DEFAULT=YES
    ;;
    *)
            # unknown option
    ;;
  esac
shift
done
readonly COUNT


#go through getting parameters via UI
start_ui

#pass parameters from user
spark_ver="${SPARK_VERSION}"
hadoop_ver="${HADOOP_VERSION}"
spark_cluster_master_url="${SPARK_MASTER_URL}"


resource_manager="${RESOURCE_MANAGER}"
yarn_hadoop_home_conf_dir="${YARN_HADOOP_HOME_CONF_DIR}"
yarn_spark_home="${YARN_SPARK_HOME}"
mesos_native_java_lib="${MESOS_NATIVE_JAVA_LIB}"
mesos_spark_executor_uri="${MESOS_SPARK_EXECUTOR_URI}"

zeppelin_port="${ZEPPELIN_PORT}"

# Main logic of installer
ZEPPELIN="zeppelin-${zeppelin_ver}"
Z_DISTR_NAME="${ZEPPELIN}-spark${spark_ver}-hadoop${hadoop_ver}.tar.gz"


download_zeppelin "${Z_DISTR_NAME}"
unpack_distr "${Z_DISTR_NAME}" "${ZEPPELIN}"

configure_interpreter "${ZEPPELIN}" \
                              "${spark_cluster_master_url}" \
                              "${yarn_spark_home}" \
                              "${mesos_spark_executor_uri}" \
                              "${yarn_hadoop_home_conf_dir}" \
                              "${mesos_native_java_lib}" \


print_instructions_to_run

persist "${spark_ver}" \
        "${hadoop_ver}" \
        "${spark_cluster_master_url}" \
        "${resource_manager}" \
        "${yarn_spark_home}" \
        "${yarn_hadoop_home_conf_dir}" \
        "${mesos_native_java_lib}" \
        "${mesos_spark_executor_uri}" \
        "${zeppelin_port}"


} # End of wrapping
