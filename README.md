# adobe-docx-conv

## Opis problemu

Aplikacja odpowiada na realny problem napotkany przy używaniu technologii frontendowej [Adobe Franklin](https://aem.live/).
Wskazana technologia pozwala używać dwóch różnych źródeł contentu, z którego tworzona jest strona internetowa:

- SharePoint (Word, Excel)
- Google Drive (Google Docs, Google Sheets)

W SharePoint pliki są przechowywane w formacie .docx, natomiast w Google Drive w formacie .gdoc. Problem pojawia się przy próbie przeniesienia plików z Google Drive do Share Point.
Google umoliwia eksport plików do formatu .docx, jednakże eksportowane pliki nie są w pełni kompatybilne z formatem .docx używanym przez Adobe Franklin na platformie Share Point.

Pierwszą niezgodnością jest oddzielanie sekcji w obu tych formatach. Adobe Franklin w Google Drive wymaga, aby sekcje były oddzielone linią poziomąd, natomiast w Share Point wymagane są 3 znaki myślnika (---). Całość jest opisana [tutaj](https://www.aem.live/docs/authoring#sections).

### Google Drive wymagany format

![Google Drive Example](./docs/googledrive.png)

### Share Point wymagany format

![Share Point Example](./docs/sharepoint.png)

### Google Drive po eksporcie do .docx - nie można tego pliku użyć w Share Point

![Google Drive Example](./docs/googledriveexport.png)

Drugim problemem jest podkreślanie linków. Przy eksporcie z Google Drive linki są podkreślone, natomiast w Share Point te podkreślenia są niepotrzebne i nierzadko sprawiają problemy.

## Opis rozwiązania

Do tej pory wszystkie pliki musiały być "naprawiane" ręcznie, natomiast ta aplikacja robi to automatycznie przez wywołanie skompilowanego programu z parametrem ze ścieką do pliku lub folderu, który ma być poprawiony.
Parametrem może być pojedynczy plik, folder lub archiwum zip.
W przypadku folderu lub archiwum zip, aplikacja przechodzi rekurencyjnie przez wszystkie pliki wewnątrz i poprawia je.
W przypadku pojedynczego pliku, aplikacja poprawia tylko ten plik.

### Przykłady użycia

Pliki, na których mona przetestować działanie aplikacji znajdują się w folderze `example`.

```bash
./conv example/starter-content
./conv example/starter-content.zip
./conv example/index.docx
```
