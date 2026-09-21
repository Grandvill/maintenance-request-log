// ==============================================================
// Jenkinsfile — CI/CD Pipeline for Maintenance Request Log
// ==============================================================

pipeline {
    agent any

    environment {
        DOCKER_BUILDKIT = '1'
        APP_ENV         = 'test'
        POSTGRES_DB     = 'maintenance_test_db'
        POSTGRES_USER   = 'test_user'
        POSTGRES_PASSWORD = 'test_password'
    }

    options {
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
        ansiColor('xterm')
    }

    stages {
        // ----------------------------------------------------------
        // Stage 1: Checkout Source Code
        // ----------------------------------------------------------
        stage('Checkout') {
            steps {
                echo 'Checking out source repository...'
                checkout scm
            }
        }

        // ----------------------------------------------------------
        // Stage 2: Backend Lint & Unit Tests
        // ----------------------------------------------------------
        stage('Backend Tests') {
            steps {
                dir('backend') {
                    echo 'Running Go static checks and tests...'
                    sh 'go vet ./...'
                    sh 'go test -v -race -coverprofile=coverage.out ./...'
                }
            }
            post {
                always {
                    echo 'Backend test stage completed.'
                }
            }
        }

        // ----------------------------------------------------------
        // Stage 3: Frontend Lint & Build Check
        // ----------------------------------------------------------
        stage('Frontend Build & Check') {
            steps {
                dir('frontend') {
                    echo 'Installing frontend dependencies and verifying Nuxt build...'
                    sh 'npm ci || npm install'
                    sh 'npm run build'
                }
            }
        }

        // ----------------------------------------------------------
        // Stage 4: Docker Compose Build Verification
        // ----------------------------------------------------------
        stage('Docker Images Build') {
            steps {
                echo 'Verifying Docker Compose build integrity...'
                sh 'cp .env.example .env'
                sh 'docker compose build'
            }
        }

        // ----------------------------------------------------------
        // Stage 5: Container Smoke Test
        // ----------------------------------------------------------
        stage('Smoke Test') {
            steps {
                echo 'Spinning up test containers to verify healthcheck...'
                sh 'docker compose up -d db backend'
                // Wait up to 30 seconds for health check endpoint
                sh '''
                    for i in $(seq 1 15); do
                        if curl -fs http://localhost:8080/health; then
                            echo "Backend healthcheck passed!"
                            exit 0
                        fi
                        echo "Waiting for service to be healthy..."
                        sleep 2
                    done
                    echo "Healthcheck timed out!"
                    exit 1
                '''
            }
            post {
                always {
                    echo 'Cleaning up test containers...'
                    sh 'docker compose down -v'
                }
            }
        }
    }

    post {
        success {
            echo '==================================================='
            echo ' Pipeline Succeeded: All tests and builds passed!   '
            echo '==================================================='
        }
        failure {
            echo '==================================================='
            echo ' Pipeline Failed: Please check the build logs above.'
            echo '==================================================='
        }
    }
}

