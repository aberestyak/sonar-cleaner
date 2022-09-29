# Sonar-cleaner

CLI tool to delete branches and analises into sonarqube's projects, which weren't analised for a log time.

> Doesn't delete projects completely, because they may have individual properties.


| Опция       | Переменная окружения    | Назначение                                                                |
| ----------- | ----------------------- | ------------------------------------------------------------------------- |
| `log-level` | SONAR_CLEANER_LOG_LEVEL | Choose log level                                                          |
| `dry-run`   | SONAR_CLEANER_DRY_RUN   | Show "outdated" projects                                                  |
| `address`   | SONARQUBE_ADDRESS       | Sonarqube address                                                         |
| `token`     | SONARQUBE_TOKEN         | Sonarqube admin token                                                     |
| `days`      | KEEP_DAYS               | Limit of days since last analysis to keep project's branches and analysis |
