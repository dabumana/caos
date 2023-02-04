# caos

[![Start with why](https://img.shields.io/badge/start%20with-why%3F-brightgreen.svg?style=flat)](https://beta.openai.com/docs/introduction/key-concepts)
[![Documentation Status](https://readthedocs.org/projects/caos-openai/badge/?version=latest)](https://caos-openai.readthedocs.io/en/latest/?badge=latest)
[![Go](https://github.com/dabumana/caos/actions/workflows/go.yml/badge.svg)](https://github.com/dabumana/caos/actions/workflows/go.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ce2f44761a6e486999eddd05b749c1be)](https://www.codacy.com/gh/dabumana/caos/dashboard?utm_source=github.com&utm_medium=referral&utm_content=dabumana/caos&utm_campaign=Badge_Grade)

### Description

Conversational assistant for openai services, includes advanced modes to modify the contextual information for specifical uses, engine, results, probabilities according to the ammount or words with predefined values for best accuracy.

### Build

Installation steps:

-   Download the following repository `git clone github.com/dabumana/caos`
-   Install dependencies:
    -   `go-gpt3`
    -   `tview`
    -   `tcell`
-   Add your API key provided from OpenAI into the `.env` file to use it with docker or export the value locally in your environment
-   Run `./clean.sh`
-   If you have Docker installed execute `./run.sh`, in any other case `./build.sh`

### Features

-   Test all the available models for **code/text/insert/similarity**

-   ##### Modes:

-   **Training mode**: Prepare your own sets based on the interaction

-   **Edit mode**: First input will be contextual the second one instructional

-   **Conversational mode**: A friendly experience with the model	

-   ##### Advanced parameters like:

    -   Temperature
    -   Topp
    -   Penalty
    -   Frequency penalty
    -   Max tokens
    -   Engine

-   Multiple results and probabilities

-   Detailed log according to UTC

### How to use?

The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.

-   OpenAI API is a set of tools and services that allow developers to create applications that use artificial intelligence (AI) technology.
-   The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.
-   To use the OpenAI API, developers must first register for an API key.

The terminal app have a conversational assistant that is designed to work with OpenAI services, able to understand natural language queries and provide accurate results,
also includes advanced modes that allow users to modify the contextual information for specific uses for example, users can adjust the engine, results, probabilities according to the amount of words used in the query, this allows for more accurate results when using longer queries.

![console.gif](docs%2Fmedia%2Fconsole.gif)

#### General parameters:

-   **Mode**: Modify the actual mode, select between (TEXT/EDIT/CODE)
-   **Engine**: Modify the model that you want to test
-   **Results**: Modify the amount of results displayed for each prompt
-   **Probabilities**: According to your setup of the temperature and topp, probably you will need to use this field to populate a more accurate response according to the possibilities of results
-   **Temperature**: If you are working with temperature, try to keep the topp in a higher values than temperature
-   **Topp**: Applies the same concept as temperature, when you are modifying this value, you need to apply a higher value for temperature
-   **Penalty**: Penalty applied to the characters an redundancy in a result completion
-   **Frequency Penalty**: Establish the frequency of the penalty threshold defined

#### Modes:

-   **Training Mode**: Prepare your fine tunning data based on the completion response, can be exported to a JSON file for further use.
-   **Edit Mode**: Use Edit mode for all the requests
    -   Press `New Conversation` and select `Edit mode` the first request will be for a completion endpoint the second based on the first request will continue editing the content in the parameters that you ask.
-   **Conversational Mode**: Use conversational AI mode request for a friendly interaction.

![details.png](docs%2Fmedia%2Fdetails.png)

### Disclaimer :warning:

This software is provided "as is" and any expressed or implied warranties, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose are disclaimed. In no event shall the author or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential.
