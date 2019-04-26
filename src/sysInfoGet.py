import psutil

if __name__ == '__main__':

    with open(r'Info/SysInfo.txt', 'w') as f:
        f.write("CPU logical num: " + str(psutil.cpu_count()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write("CPU physical cores num: " +
                str(psutil.cpu_count(logical=False)) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.cpu_times()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write("CPU using in near 10 sec:\n")
    for x in range(1):
        with open(r'Info/SysInfo.txt', 'a') as f:
            f.write(str(psutil.cpu_percent(interval=1, percpu=True)) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.virtual_memory()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.swap_memory()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.disk_partitions()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.disk_usage('/')) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.disk_io_counters()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.net_io_counters()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.net_if_addrs()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.net_if_stats()) + "\n")
    with open(r'Info/SysInfo.txt', 'a') as f:
        f.write(str(psutil.net_connections()) + "\n")
    