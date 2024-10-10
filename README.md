# YouTube Learning Mode - Video Processing Service

This is the Golang-based **Video Processing Service** for the YouTube Learning Mode project. This service captures video snapshots at specified timestamps, stores them in a Redis database, and aims to enhance AI responses by providing visual context from video frames. It is an integral part of a microservices architecture that includes the main backend, AI service, and a Python-based video information extraction service.

## Features

- **Snapshot Capture**: Downloads YouTube videos using `yt-dlp` and captures snapshots at specified timestamps using `ffmpeg`.
- **Redis Integration**: Stores captured snapshots along with metadata in Redis, allowing for easy retrieval by other services.
- **Microservices Architecture**: Designed to work alongside other services like the main backend, AI service, and YouTube Info Service.
- **REST API**: Provides a REST API endpoint for processing snapshot requests.

## Prerequisites

- **Go**: Version 1.16 or higher.
- **Redis**: Used for storing snapshot metadata.
- **yt-dlp**: For downloading YouTube videos.
- **ffmpeg**: For extracting frames from videos.
- **Docker**: For containerizing the service.
- **Docker Compose**: For orchestrating all microservices in a shared network.

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/AnasKhan0607/Youtube-Learning-Mode-Video-Processing-Service.git
cd Youtube-Learning-Mode-Video-Processing-Service
```

### 2. Set Up Environment

> **Note:** This service doesn't require a `.env` file but expects Redis to be running at `redis:6379` when containerized. Adjust the Redis connection if needed in the `initRedis` function.

### 3. Build and Run Using Docker

To build the Docker image and run the service, use:

```bash
docker build -t video-processing-service .
docker run -p 8081:8081 --network="host" video-processing-service
```

Alternatively, if using **Docker Compose** with other services:

1. Make sure all services are in the same directory with a `docker-compose.yml` file.
2. Run:

```bash
docker-compose up --build
```

## API Endpoints

### 1. /process-snapshots (POST)

Processes a YouTube video URL and timestamps, captures snapshots at those timestamps, and stores them in Redis.

- **Request Body**:

  ```json
  {
    "video_id": "VIDEO_ID",
    "video_url": "https://www.youtube.com/watch?v=VIDEO_ID",
    "timestamps": ["00:01:00", "00:02:00", "00:03:00"]
  }
  ```

- **Response**:

  ```json
  {
    "message": "Snapshots stored successfully"
  }
  ```

- **Description**: The request body includes the `video_id`, `video_url`, and a list of `timestamps`. The service downloads the video using `yt-dlp`, extracts frames at the specified timestamps using `ffmpeg`, and stores the snapshot paths in Redis under the provided `video_id`.

## Project Structure

```
├── cmd/
│   └── main.go               # Entry point of the application
├── Dockerfile                # Docker configuration for containerizing the service
├── go.mod                    # Go module dependencies
├── go.sum                    # Checksums for Go modules
└── README.md                 # Project documentation
```

## Dependencies

- **Redis Go Client**: For connecting to the Redis database.
- **yt-dlp**: Downloads videos from YouTube.
- **ffmpeg**: Captures frames from videos.
- **Go HTTP Package**: For handling REST API requests.

## How It Works

1. **Download Video**: The service uses `yt-dlp` to download a video from a given URL.
2. **Capture Snapshot**: For each specified timestamp, `ffmpeg` is used to extract a frame from the downloaded video.
3. **Store in Redis**: Each captured snapshot is stored in Redis as a JSON object, associated with the video ID and timestamp.
4. **Response**: After processing all timestamps, a success message is returned to the client.

## Future Enhancements

- **Integrate AI Feedback**: Use the stored snapshots to provide visual cues to the AI service, helping improve response accuracy.
- **Optimize Video Downloading**: Implement mechanisms to avoid redundant downloads if a video has been processed previously.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.

2. Create a new branch:

   ```bash
   git checkout -b feature/YourFeature
   ```

3. Commit your changes:

   ```bash
   git commit -am 'Add some feature'
   ```

4. Push to the branch:

   ```bash
   git push origin feature/YourFeature
   ```

5. Open a Pull Request.


### Important Notes:

- This service relies heavily on **yt-dlp** and **ffmpeg** for video processing. Ensure these tools are installed on your system or included in the Docker image.
- It is designed to work as part of a larger microservices architecture. Running it standalone will require manual handling of Redis and other dependent services.
