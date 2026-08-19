 cd docker-infra

 docker compose --env-file .env.example down -v

 echo $env:NOTIFICATION_DB_DSN

 $env:DB_REPOS_DIR=".."

 docker compose --env-file .env.example up postgres -d

 docker compose --env-file .env.example ps

 docker compose --env-file .env.example --profile tooling run --rm liquibase-notification update



----------------------------------------------------------------

# Conexión al servidor PostgreSQL desde pgAdmin

1. Abrir **pgAdmin**.

2. En el panel izquierdo:
   - Clic derecho en **Servers**.
   - Seleccionar **Register → Server...**

3. En la pestaña **General**:
   - **Name:** `design-software-develop`

4. En la pestaña **Connection** ingresar:

   - **Host name/address:** `localhost`
   - **Port:** `15432`
   - **Maintenance database:** `design-software-develop`
   - **Username:** `design_software_user`
   - **Password:** `change-me`

5. Activar **Save password** si se desea guardar la contraseña.

6. Presionar **Save**.

7. Si la conexión es correcta, aparecerá el servidor en pgAdmin y se podrá acceder a:

   **Servers → design-software-develop → Databases → design-software-develop**
