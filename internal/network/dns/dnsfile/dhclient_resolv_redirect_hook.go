package dnsfile

const DhclientResolvRedirectHook = `
#!/bin/sh

# This file is managed by edge-guard, please don't modify.

# update /etc/resolv.conf based on received values
make_resolv_conf() {
    local new_resolv_conf

    echo "====EDGE-GUARD DHCP START===="

    # DHCPv4
    if [ -n "$new_domain_search" ] || [ -n "$new_domain_name" ] ||
       [ -n "$new_domain_name_servers" ]; then
        resolv_conf=$(readlink -f "/etc/resolv.conf" 2>/dev/null) ||
            resolv_conf="/etc/resolv.conf"

        mkdir -p /edge/resolv.d
        new_resolv_conf="/edge/resolv.d/dhclient.${interface}"
        wait_for_rw "$new_resolv_conf"
        rm -f $new_resolv_conf

        if [ -n "$new_domain_name" ]; then
            echo domain ${new_domain_name%% *} >>$new_resolv_conf
        fi

        if [ -n "$new_domain_search" ]; then
            if [ -n "$new_domain_name" ]; then
                domain_in_search_list=""
                for domain in $new_domain_search; do
                    if [ "$domain" = "${new_domain_name}" ] ||
                       [ "$domain" = "${new_domain_name}." ]; then
                        domain_in_search_list="Yes"
                    fi
                done
                if [ -z "$domain_in_search_list" ]; then
                    new_domain_search="$new_domain_name $new_domain_search"
                fi
            fi
            echo "search ${new_domain_search}" >> $new_resolv_conf
        elif [ -n "$new_domain_name" ]; then
            echo "search ${new_domain_name}" >> $new_resolv_conf
        fi

        if [ -n "$new_domain_name_servers" ]; then
            for nameserver in $new_domain_name_servers; do
                echo nameserver $nameserver >>$new_resolv_conf
            done
        else # keep 'old' nameservers
            sed -n /^\w*[Nn][Aa][Mm][Ee][Ss][Ee][Rr][Vv][Ee][Rr]/p $resolv_conf >>$new_resolv_conf
        fi

        if [ -f $resolv_conf ]; then
            chown --reference=$resolv_conf $new_resolv_conf
            chmod --reference=$resolv_conf $new_resolv_conf
        fi

    echo "====EDGE-GUARD DHCP END===="
    # DHCPv6 ignored
    fi
}`
