import psutil
import sys
import io


if __name__ == '__main__':
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf8')
    f = open("test.txt", 'w')
    old = sys.stdout
    sys.stdout = f
    print(psutil.test().decode('gbk'))
    sys.stdout = old
    f.close()
