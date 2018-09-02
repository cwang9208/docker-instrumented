#!/bin/bash

vm_nr=1000
core_nr=3
core_ini=1
#iter=1000000
iter=2800000

source /root/persistent/normal-distribution.input

create_vm_json()
{
    gid=${1}
    cpuid=${2}
    ip=${3}
    mac=${4}
    fname=${5}

    cat > ${fname} <<EOF
{
    "name" : "awslambda.${gid}",

    "kernel" : "/root/awslambda_xs",
    "cmdline" : "${ip} ${iter}",

    "memory" : 16,
    "vcpus" : {
        "count" : 1,
        "cpumap" : [[${cpuid}]]
    },

    "vifs" : [
        {
            "mac" : "${mac}",
            "ip" : "${ip}",
            "bridge" : "xenbr"
        }
    ],

    "xen" : {
        "dev_method" : "xenstore"
    }
}
EOF
}

create_vm_domconf()
{
    gid=${1}
    cpuid=${2}
    ip=${3}
    mac=${4}
    fname=${5}

    cat > ${fname} <<EOF
kernel = "/root/awslambda_xs"
memory = "16"
cpus = ["0"]
vcpus = "1"
name = "aws.${gid}"
cmdline = "${ip} ${iter}"
vif  = [ 'mac=${mac},bridge=xenbr,ip=${ip}' ]
on_poweroff = "destroy"
on_crash = "preserve"
on_reboot = "preserve"
EOF
}

for i in $(seq 0 $(( ${vm_nr} - 1 ))) ; do
    cpu=$(( (${i} % ${core_nr}) + ${core_ini} ))

    ipthird=$(( ${i} / 254 ))
    ipfourth=$(( (${i} % 254) + 1 ))
    ip=10.128.${ipthird}.${ipfourth}
    macthird=$(printf "%02x" ${ipthird})
    macfourth=$(printf "%02x" ${ipfourth})
    mac=16:27:29:a3:${macthird}:${macfourth}

    create_vm_json ${i} ${cpu} ${ip} ${mac} /tmp/aws.json
    chaos --no-noxs create /tmp/aws.json 2>>create.log &
    #create_vm_domconf ${i} ${cpu} ${ip} ${mac} /tmp/aws.xen
    #xl create /tmp/aws.xen

    #sleep ${interarrival[$i]}
    j=$(echo "${interarrival[$i]} * 10" | bc)
    sleep $j
done

