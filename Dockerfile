FROM golang:1.19-buster as builder

RUN printf "deb https://mirrors.aliyun.com/debian/ buster main non-free contrib\n\
deb https://mirrors.aliyun.com/debian-security buster/updates main\n\
deb https://mirrors.aliyun.com/debian/ buster-updates main non-free contrib\n\
deb https://mirrors.aliyun.com/debian/ buster-backports main non-free contrib\n\
" > /etc/apt/sources.list
RUN apt update 
RUN apt-get -y install zip unzip
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /financial_statement
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN unzip -d /usr/local/go/src/ dm-go-driver.zip

# RUN wget https://github.com/bblanchon/pdfium-binaries/releases/latest/download/pdfium-linux-x64.tgz
ENV PDFIUM_PATH /opt/pdfium
RUN mkdir -p ${PDFIUM_PATH} && tar xzvf pdfium-linux-x64.tgz -C ${PDFIUM_PATH}
RUN mkdir -p /usr/lib/pkgconfig && printf "prefix=${PDFIUM_PATH}\n\
libdir=${PDFIUM_PATH}/lib\n\
includedir=${PDFIUM_PATH}/include\n\
\n\
Name: PDFium\n\
Description: PDFium\n\
Version: 5254\n\
Requires:\n\
\n\
Libs: -L${PDFIUM_PATH}/lib -lpdfium\n\
Cflags: -I${PDFIUM_PATH}/include" > /usr/lib/pkgconfig/pdfium.pc

RUN export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$PDFIUM_PATH/lib 
RUN go build -ldflags="-s -w" -o /financial_statement/financial_statement cmd/api-server/apiserver.go



# libreoffice安装太慢了，构建了一个base镜像提高镜像构建速度，libreoffice:1.7.4.3_rc2_0ubuntu0.20.04.1_lo1镜像Dockerfile详见 Dockerfile-libreoffice:1.7.4.3 文件
FROM registry.intsig.net/textin_global/libreoffice:1.7.4.3_rc2_0ubuntu0.20.04.1_lo1
ENV LD_LIBRARY_PATH /opt/pdfium/lib
ENV TZ Asia/Shanghai
RUN rm -f /financial_statement    
COPY --from=builder /opt/pdfium /opt/pdfium
WORKDIR /app
COPY --from=builder /financial_statement/financial_statement /app/financial_statement
COPY --from=builder /financial_statement/configs /app/configs
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./financial_statement", "-c", "configs/api-server-debug.yaml"]
