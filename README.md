# Galore Services
![Header](./.github/header.png)
Repository holding all the services required for galore to work.

Currently, it has the following services:
- API-Service(main): all the main api's
- Image Generation: service for generating images
- Embedding Generation: service for generating embeddings
- Categorizer: service for categorizing cocktails based on cocktail properties

\*The project is a rewrite, initially it was written using Supabase & Typescript.

The following repository is a part of a larger project named Galore: 
- [**Check out the Android App**](https://github.com/m1thrandir225/galore-android)
- [**Check out the iOS App**](https://github.com/m1thrandir225/galore-ios)
- [**Check out the Admin Dashboard(WIP)**](https://github.com/m1thrandir225/galore-dashboard)
- [**Check out the Landing Site**](https://github.com/m1thrandir225/galore-landing)

# Run Locally 
Each of the services contains a `.env.example` file that shows what parameters you need to run each 
service separately. 

If you want to run them using a `docker-compose` setup there are multiple also in the base directory 
`.env.example` files for each service that requires some kind of configuration.

---
License MIT © [Sebastijan Zindl](./LICENSE)
