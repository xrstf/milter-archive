# milter Archive

This repository contains a simple Go application that functions as a sendmail milter
that passively reads incoming e-mails and writes them into individual `.eml` files.

It's useful as a quick&dirty debugging solution for incoming messages or be used to
retain training data for spam filters when messages would otherwise be redirected
to external servers. For this usecase the `-spam` flag can be used to write e-mails
into `ham` or `spam` directory based on the presence of an `X-Spam` header (most
likely inserted by a previous milter, like rspamd).

A bunch of code is based on https://github.com/phalaaxx/pf-milters.

## Usage

    $ go get github.com/xrstf/milter-archive
    $ ./milter-archive -spam -target /var/lib/archive
    2019/05/04 21:01:28 Listening on tcp:127.0.0.1:47256 ...
    2019/05/04 21:01:28 Accepting connections now.

After receiving a bunch of e-mails, `/var/lib/archive` might look like this:

    $ ls -lah /var/lib/archive/spam
    drwxr-xr-x 2 root root 4.0K May  4 18:58 .
    drwx------ 3 root root 4.0K May  4 18:48 ..
    -rw-r--r-- 1 root root 2.7K May  4 18:31 2019-05-04T183102.865Z-re-reliable-supplier-for-led-sign-and-display.eml
    -rw-r--r-- 1 root root 4.1K May  4 18:35 2019-05-04T183557.810Z-re-iron-core-metal-crawler.eml
    -rw-r--r-- 1 root root 2.6K May  4 18:37 2019-05-04T183715.052Z-re-watch.eml
    -rw-r--r-- 1 root root 2.6K May  4 18:48 2019-05-04T184804.922Z-re-2019-newest-fashion-ladies-handbag.eml

## License

MIT
