# Dark Throne Automate

Dark Throne Automate is a project designed to automate tasks for the Dark Throne game. This tool aims to simplify repetitive actions, improve efficiency, and enhance the overall gaming experience.

## Features

- Automates common in-game tasks.
- Written in Go for better performance and maintainability.
- Lightweight and easy to configure.

## Requirements

- Go (version 1.24.3 or higher)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/rihoj/darkThroneAutomate.git
   ```
2. Navigate to the project directory:
   ```bash
   cd darkThroneAutomate
   ```
3. Build the project:
   ```bash
   go build -o darkThroneAutomate
   ```

## Usage

1. Configure the `.env` file with necessary parameters (e.g., API keys, thresholds).
2. Run the executable with the required arguments:
   ```bash
   ./darkThroneAutomate <should_attack> <start_page> <end_page> <gold_threshold>
   ```
   - `should_attack`: Set to `true` to enable attacks, or `false` to only evaluate targets.
   - `start_page`: the starting pagination page
   - `end_page`: the ending pagination page
   - `gold_threshold`: Minimum gold amount to consider a target.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes and push the branch:
   ```bash
   git commit -m "Description of changes"
   git push origin feature-name
   ```
4. Open a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is intended for educational purposes only. Use it responsibly and ensure compliance with the game's terms of service.
