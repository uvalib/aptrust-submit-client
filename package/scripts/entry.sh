#
#
#

# run application
cd bin; ./aptrust-submit-client \
   -port      ${SERVICE_PORT} \
   -jwtkey    ${JWT_KEY} \
   -group     ${AUTH_GROUP} \
   -busname   ${EVENT_BUS_NAME} \
   -eventsrc  ${EVENT_SRC_NAME} \
   -dbhost    ${DB_HOST} \
   -dbport    ${DB_PORT} \
   -dbname    ${DB_NAME} \
   -dbuser    ${DB_USER} \
   -dbpass    ${DB_PASSWORD}

#
# end of file
#
