FROM alpine:3.20

# Define the variable without a default value
ARG APP_NAME

# Check if APP_NAME was passed; fail early if empty
RUN if [ -z "$APP_NAME" ]; then echo "Error: APP_NAME buildarg is required." && exit 1; fi

# Create a system user named dynamically after APP_NAME
RUN adduser -D -H -s /sbin/nologin ${APP_NAME}

WORKDIR /app

# Copy files and assign ownership dynamically
COPY --chown=${APP_NAME}:${APP_NAME} build/${APP_NAME}.linux /usr/local/bin/${APP_NAME}
COPY --chown=${APP_NAME}:${APP_NAME} build/${APP_NAME}.zip /app/

# Switch to the dynamic user
USER ${APP_NAME}

EXPOSE 8080

# Convert build-time ARG to runtime ENV for CMD expansion
ENV RUN_APP=${APP_NAME}
CMD ["sh", "-c", "exec ${RUN_APP}"]
