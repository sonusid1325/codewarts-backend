FROM ubuntu:22.04

# Avoid prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install core utilities and common shell tools
RUN apt-get update && apt-get install -y \
    curl \
    vim \
    nano \
    grep \
    findutils \
    coreutils \
    iputils-ping \
    net-tools \
    sudo \
    file \
    procps \
    psmisc \
    hostname \
    tar \
    gzip \
    cron \
    iproute2 \
    dnsutils \
    traceroute \
    wget \
    sed \
    gawk \
    openssh-client \
    rsync \
    && rm -rf /var/lib/apt/lists/*

# Create player user
RUN useradd -m -s /bin/bash player && \
    echo "player:player" | chpasswd && \
    usermod -aG sudo player

# Setup passwordless sudo for the player (helps if they want to explore, but keeps it safe inside the container)
RUN echo "player ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Create game directory structures
RUN mkdir -p /var/mainframe/archives \
    && mkdir -p /var/mainframe/cores \
    && mkdir -p /var/mainframe/terminal_history \
    && echo "Archive systems online." > /var/mainframe/archives/system_backup.txt \
    && echo "Core offline." > /var/mainframe/cores/core_state.txt

# Pre-populate logs for Chapter 2
RUN echo "[2026-05-26 00:00:01] Mainframe system starting up...\n[2026-05-26 00:00:15] Security modules initialized.\n[2026-05-26 00:00:30] WARNING: Breach detected in Sector 7." > /var/log/mainframe.log

# Pre-populate logs for Chapter 4
RUN echo "[2026-05-26 00:01:00] SYSTEM BOOTING: SUCCESS\n\
[2026-05-26 00:02:15] SECTOR 1: ACTIVE (TEMP 45C)\n\
[2026-05-26 00:03:30] SECTOR 2: ACTIVE (TEMP 48C)\n\
[2026-05-26 00:04:12] WARNING: FUEL LEVEL LOW\n\
[2026-05-26 00:05:00] SECTOR 3: ACTIVE (TEMP 52C)\n\
[2026-05-26 00:05:10] CRITICAL: REACTOR OVERHEAT PASSWORD: CORE_TEMP_CRITICAL_9982\n\
[2026-05-26 00:06:00] SECTOR 4: ACTIVE (TEMP 55C)\n\
[2026-05-26 00:07:05] ALERT: COOLING SYSTEM FAILURE\n\
[2026-05-26 00:08:00] SECTOR 5: INACTIVE (MAINTENANCE)" > /var/log/reactor.log

# Setup player home directory and Chapter 1 secret keys
USER player
WORKDIR /home/player

RUN echo "GRID_PASS_99" > /home/player/.cyber_key && \
    chmod 600 /home/player/.cyber_key

# Set default command to bash
CMD ["/bin/bash"]
