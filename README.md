# Catuaba Framework

Fast, super productive and aphrodisiac web-framework written in golang

![Catuaba new command example](catuaba-new-command.gif)

# Setup

Install using go get command
```
go get -u github.com/dayvsonlima/catuaba
```

## On Windows
Install the [tdm-gcc](https://jmeubank.github.io/tdm-gcc/) 

# Command List
| Command | Description | Usage |
| --- | --- | --- |
|new | Create a New project | `catuaba new <project-name>`
|server| Starts the catuaba web server | Run `catuaba server` in the project root
|generator| Call some catuaba generator | `catuaba generator <generator-name> ...<generator-params>`

# Generators List
| Command | Description | Usage |
| --- | --- | --- |
|scaffold| scaffold is a full set of model, controller and router for one RESTFUL resource | `catuaba g scaffold <model-name> ...<attribute:type>`
|model| generates a new model file in `app/models` directory | `catuaba g model <model-name> ...<attribute:type>`
|controller| generates a new controller package in `app/controller` | `catuaba g controller <controller-name> ...<method-name (index, new, create, update, delete, all)>`. if you don't specify the method name, catuaba will create them all by default

TO DO
 - install a rake-like tool
  https://github.com/magefile/mage

 - Migrations structure like this
  https://www.youtube.com/watch?v=9QJ5mQbRUPU

 - https://github.com/jinzhu/inflection

