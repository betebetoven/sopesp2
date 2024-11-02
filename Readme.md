# Manual Técnico de Proyecto en Kubernetes

Este manual detalla los comandos y configuraciones necesarias para desplegar y gestionar los componentes del proyecto, incluyendo `Kubernetes`, `Kafka`, `Redis`, `Grafana`, `Prometheus` y `Locust`. El proyecto utiliza diferentes servicios para procesar datos de estudiantes, distribuir tráfico según facultad y disciplina, y monitorear el estado del sistema en tiempo real.

## Estructura del Proyecto

La estructura de directorios del proyecto es la siguiente:

- **`agronomia`**: Servicio en Go que procesa datos de estudiantes para la facultad de Agronomía.
- **`agronomia-monitor`**: Monitoreo de eventos de Agronomía con Redis y Pub/Sub.
- **`deployment`**: Archivos de despliegue en Kubernetes.
- **`deployment_dis2`**: Despliegues específicos para los servicios de disciplinas.
- **`deployment_ing`**: Despliegues para los servicios de Ingeniería.
- **`grafana-data`**: Configuración de dashboards para Grafana.
- **`ingenieria`**: Servicio de la facultad de Ingeniería implementado en Rust.
- **`ingenieria-monitor`**: Monitoreo de eventos de Ingeniería.
- **`kafka-results`**: Servicio que envía resultados a Kafka.
- **`locust`**: Generador de tráfico utilizando Locust.
- **`monitoring`**: Configuraciones para Prometheus y monitoreo.

### Descripción de los Archivos Clave y Comandos Utilizados

1. **Despliegue de Servicios para Facultades y Disciplinas**  
   Cada servicio para facultad y disciplina se implementa mediante `Kubernetes` con un `Deployment` y un `Service` para la comunicación vía gRPC.

    ```yaml
    # Ejemplo: Despliegue Kubernetes para el servicio de Agronomía
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: agronomia-deployment
    spec:
      replicas: 2
      selector:
        matchLabels:
          app: agronomia-service
      template:
        metadata:
          labels:
            app: agronomia-service
        spec:
          containers:
          - name: agronomia
            image: betebetoven/agronomia:latest
            ports:
            - containerPort: 8081
    ```

    Para cada servicio de disciplina (`discipline1`, `discipline2`, `discipline3`), se configura un despliegue en Kubernetes de manera similar, cambiando el nombre del contenedor y puerto, como se ve en los archivos `discipline1-deployment.yaml`, `discipline2-deployment.yaml`, y `discipline3-deployment.yaml`.

2. **Configuración de Ingress para Redirección de Tráfico**  
   `Ingress` permite redirigir el tráfico según el endpoint de facultad en las peticiones HTTP. Aquí se redirige el tráfico a `/ingenieria` o `/agronomia`.

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: faculty-ingress
    spec:
      rules:
      - host: "example.com"
        http:
          paths:
          - path: /ingenieria
            pathType: Prefix
            backend:
              service:
                name: ingenieria-service
                port:
                  number: 8080
          - path: /agronomia
            pathType: Prefix
            backend:
              service:
                name: agronomia-service
                port:
                  number: 8081
    ```

3. **Configuración de Autoescalado de Pods (HPA)**  
   Configuramos un `Horizontal Pod Autoscaler` para escalar automáticamente los pods en función del uso de CPU:

    ```bash
    kubectl autoscale deployment agronomia-deployment --cpu-percent=50 --min=1 --max=5
    ```

4. **Simulación de Tráfico con Locust**  
   `Locust` se utiliza para simular tráfico hacia los servicios de las facultades. Genera tráfico dirigido a los endpoints `/ingenieria` y `/agronomia`.

    ```python
    # locustfile.py
    from locust import HttpUser, task, between
    import json, random

    class StudentUser(HttpUser):
        wait_time = between(1, 3)

        @task(1)
        def health_check_ingenieria(self):
            self.client.get("/ingenieria/health")

        @task(1)
        def health_check_agronomia(self):
            self.client.get("/agronomia/health")

        @task(3)
        def add_ingenieria_student(self):
            test_case = {
                "student": "Estudiante Ejemplo",
                "age": 20,
                "faculty": "Ingenieria",
                "discipline": 1
            }
            self.client.post("/ingenieria/add_student", json=test_case)
    ```

5. **Configuración de Kafka para Almacenar Resultados**  
   Los servicios de cada disciplina envían los resultados a un sistema de colas `Kafka`, almacenándolos en los tópicos `winners` y `losers`.

    ```go
    // Ejemplo de envío a Kafka en main.go
    msg := &sarama.ProducerMessage{
        Topic: "student-results",
        Value: sarama.StringEncoder(result),
    }
    ```

6. **Uso de Redis y Configuración de Grafana**  
   Redis se utiliza para almacenar los resultados de las disciplinas, mientras que Grafana muestra visualizaciones en tiempo real.

    ```yaml
    # Configuración de Prometheus para Monitoreo
    apiVersion: monitoring.coreos.com/v1
    kind: PrometheusRule
    metadata:
      name: prometheus-rules
    spec:
      groups:
      - name: monitoring-rules
        rules:
        - alert: HighPodCPU
          expr: sum(rate(container_cpu_usage_seconds_total[5m])) by (pod) > 0.5
          for: 5m
          labels:
            severity: "critical"
          annotations:
            summary: "Uso alto de CPU"
            description: "Pod {{ $labels.pod }} usa demasiada CPU"
    ```

7. **Configuración y Visualización en Grafana**  
   Los datos de Redis se muestran en dashboards de Grafana para visualizar el estado de los estudiantes por facultad, disciplina, y resultados. También se muestran métricas de Prometheus sobre el estado del clúster y el rendimiento del sistema.

    ```yaml
    # Despliegue de Grafana en Kubernetes
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: grafana
    spec:
      replicas: 1
      template:
        spec:
          containers:
          - name: grafana
            image: grafana/grafana:latest
            ports:
            - containerPort: 3000
    ```

8. **Redis para Almacenamiento de Resultados**

    ```go
    // Almacenamiento de resultados en Redis
    key := fmt.Sprintf("faculty:%s:discipline:%d:result:%s", student.Faculty, student.Discipline, result)
    rdb.Incr(ctx, key)
    ```

### Comandos de Despliegue y Monitoreo

1. **Construcción y Publicación de Imágenes Docker**

    ```bash
    docker build -t betebetoven/agronomia .
    docker push betebetoven/agronomia
    ```

2. **Despliegue en Kubernetes**

    ```bash
    kubectl apply -f agronomia-deployment.yaml
    kubectl apply -f agronomia-service.yaml
    ```

3. **Escalado de Despliegues**

    ```bash
    kubectl scale deployment agronomia-deployment --replicas=3
    ```

4. **Configuración del Autoscalador de Pods (HPA)**

    ```bash
    kubectl autoscale deployment agronomia-deployment --cpu-percent=60 --min=2 --max=10
    ```

5. **Despliegue de Redis, Grafana y Prometheus con Helm**

    ```bash
    # Instalar Redis
    helm install redis bitnami/redis

    # Instalar Grafana
    helm install grafana stable/grafana

    # Instalar Prometheus
    helm install prometheus prometheus-community/prometheus
    ```

6. **Verificación del Estado de Ingress**

    ```bash
    kubectl get ingress
    ```

7. **Verificación de Servicios y Pods**

    ```bash
    kubectl get services
    kubectl get pods
    ```

### Resumen

Este manual técnico proporciona una guía detallada de los comandos y configuraciones para desplegar la arquitectura del proyecto. La estructura incluye la gestión de tráfico, almacenamiento y visualización de datos en tiempo real, y el escalado automático de servicios según la demanda. Redis almacena la información de estudiantes, Grafana y Prometheus ofrecen monitoreo y visualización, mientras que Kafka distribuye los resultados de disciplinas ganadores y perdedores.