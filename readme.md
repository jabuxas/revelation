# revelation

this app triggers the dbus filepicker to select a file, then uploads it to a pastebin service.

## how it works

- when you run the app, it opens the system filepicker using dbus.

- after selecting a file, the app uploads it to the pastebin service.

- the app returns the URL of the uploaded file.

## requirements

- dbus (for the filepicker)

- a pastebin service like [0x0.st](https://0x0.st)

## usage

1. set the AUTH_KEY and AUTH_PARAM environment variable if your pastebin service requires an API key.

1. run the app:

   ```bash
   go run .
   ```

1. select a file using the filepicker.

1. the app will copy the url of the uploaded file to clipboard

## notes

- this app is designed for linux systems with dbus support.
- make sure the pastebin service is accessible and the API key is valid.
