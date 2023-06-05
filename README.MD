# Go Playground Countdowns
GoPlayground Countdowns is a simple project that is my learning playground for exploring the capabilities and features of programming language Go, particularly in creating dynamic GIF countdowns for email campaigns. The project uses several Go libraries and technologies to generate, manipulate, and animate images.

This project is essentially a GIF generator that creates countdown timers suitable for email campaigns. The generated GIF has 60 frames, each representing one minute of an hour. When a request to generate a GIF is received, new gifs are created automatically for next 60 minutes, creating images on the fly to reduce response times.

Each frame shows the countdown timer with the remaining time, with one minute being subtracted for each subsequent frame. This gives the illusion of a ticking timer when the GIF is played.

Firstly, placeholder image is generated which has little bit lower quality to reduce response time, meanwhile new GIFs are created with better quality and cached to improve both quality and response times for future requests.

Features
- **GIF Countdown**: Generate dynamic countdown GIFs, useful for visual countdowns for events, sales, promotions, etc., in email campaigns.
- **Caching**: Implementing efficient caching strategies using Go variables to improve the performance of countdown generation adn response times.
  

## Getting Started
### Prerequisites
- Go 1.13 or later.
- Basic knowledge of Go programming.

### Installation
To get a local copy up and running, follow these simple steps:

1. Clone the repo:
 ```shell 
 git clone https://github.com/bxxf/goplayground-countdowns.git
 ```

2. Navigate to the project directory:
 ```shell 
 cd goplayground-countdown
 ```

3. Run the project:
 ```shell 
 go run main.go
 ```

You can find demo [here](https://countdown-gifs-test.fly.dev/countdown?date=2023-06-15)