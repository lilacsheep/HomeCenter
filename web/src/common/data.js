export default {
    data : {
        dnsAddrs: [
            { value: "8.8.8.8", label: "Google" },
            { value: "208.67.222.222", label: "OpenDNS" },
            { value: "114.114.114.114", label: "电信" },
            { value: "1.1.1.1", label: "Cloudflare" },
            { value: "223.5.5.5", label: "阿里" },
            { value: "180.76.76.76", label: "百度" },
            { value: "119.29.29.29", label: "DNSPod" },
            { value: "202.141.162.123", label: "中科大防污染DNS(电信)" },
            { value: "202.141.176.93", label: "中科大防污染DNS(移动)" },
        ],
    },
    set_data: function(value) {
        this.data
    }
}