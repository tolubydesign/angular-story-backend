FROM postgres
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB world
COPY world.sql /docker-entrypoint-initdb.d/

# Run 
# $ docker build -t my-postgres-db ./
# $ docker images -a
# $ docker run -d --name my-postgresdb-container -p 5432:5432 my-postgres-db
