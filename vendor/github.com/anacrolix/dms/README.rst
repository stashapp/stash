dms
===

.. image:: https://circleci.com/gh/anacrolix/dms.svg?style=svg
    :target: https://circleci.com/gh/anacrolix/dms

dms is a UPnP DLNA Digital Media Server. It runs from the terminal, and serves
content directly from the filesystem from the working directory, or the path
given. The SSDP component will broadcast and respond to requests on all
available network interfaces.

dms advertises and serves the raw files, in addition to alternate transcoded
streams when it's able, such as mpeg2 PAL-DVD and WebM for the Chromecast. It
will also provide thumbnails where possible.

dms uses ``ffprobe``/``avprobe`` to get media data such as bitrate and duration, ``ffmpeg``/``avconv`` for video transoding, and ``ffmpegthumbnailer`` for generating thumbnails when browsing. These commands must be in the ``PATH`` given to ``dms`` or the features requiring them will be disabled.

.. image:: https://lh3.googleusercontent.com/-z-zh7AzObGo/UEiWni1cQPI/AAAAAAAAASI/DRw9IoMMiNs/w497-h373/2012%2B-%2B1

Installing
==========

Assuming ``$GOPATH`` and Go have been configured already::

    $ go get github.com/anacrolix/dms

Ensure ``ffmpeg``/``avconv`` and/or ``ffmpegthumbnailer`` are in the ``PATH`` if the features depending on them are desired.

To run::

    $ "$GOPATH"/bin/dms

Running DMS as a systemd service
=================================

A sample systemd `.service` file has been `provided <helpers/systemd/dms.service>`_ to assist in running DMS as a system service.

Running DMS as a FreeBSD service
================================

Install the `provided <helpers/bsd/dms>`_ service file to /etc/rc.d or /usr/local/etc/rc.d
add ``dms_enable="YES"``, and optionally ``dms_root="/path/to/my/media"`` and ``dms_user="myuser"`` to your /etc/rc.conf

Known Compatible Players and Renderers
======================================

 * Probably all Panasonic Viera TVs.
 * Android's BubbleUPnP and AirWire
 * Chromecast
 * VLC
 * LG Smart TVs, with varying success.
