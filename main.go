package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type AgentState struct {
	Step        int
	TargetMet   bool
	CurrentData string
}

type Agent struct {
	State AgentState
}

// 🛠️ FUNCIÓN M_C_P SIMULADA: Nuestra herramienta de conocimiento externo
// Devuelve documentación técnica de la API según lo que pida la IA
func (a *Agent) CallKnowledgeServer(query string) string {
	fmt.Printf("🔍 [Servidor MCP Local] Buscando documentación para: '%s'...\n", query)
	
	// Simulamos la base de datos de conocimiento de Google Developer
	if strings.Contains(strings.ToLower(query), "temperatura") || strings.Contains(strings.ToLower(query), "gpu") {
		return "MANUAL DE HARDWARE: Las GPUs estables operan por debajo de 85°C. A partir de 95°C se debe suspender el minero XMRig para evitar degradación de silicio."
	}
	return "CONOCIMIENTO: El sistema opera en modo normal sin alertas registradas."
}

func (a *Agent) Perceive() {
	a.State.Step++
	fmt.Printf("\n[Paso %d] 👁️ Escaneando sensores del sistema local...\n", a.State.Step)
	a.State.CurrentData = "TELEMETRÍA ACTUAL: GPU detectada a 97°C."
}

func (a *Agent) Think() string {
	fmt.Println("🧠 Consultando al cerebro del agente en Google...")

	apiKey := "TU_CLAVE_API_AQUÍ"
   
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=" + apiKey

	// Aquí aplicamos Prompt Engineering avanzado (lo que enseña la Unidad 2)
	// Le damos a la IA la capacidad de usar nuestra herramienta MCP si le hace falta información
	prompt := fmt.Sprintf(`Eres un agente inteligente. Datos actuales: %s.
Reglas obligatorias:
1. Si necesitas saber las acciones de seguridad para mitigar los datos actuales, responde exactamente: SOLICITAR_CONOCIMIENTO_HARDWARE
2. Si ya conoces la acción técnica a tomar basándote en los datos, responde únicamente: EJECUTAR_APAGADO_SEGURIDAD
No añadas nada más.`, a.State.CurrentData)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "ERROR"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var responseData struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	json.Unmarshal(body, &responseData)

	if len(responseData.Candidates) > 0 && len(responseData.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(responseData.Candidates[0].Content.Parts[0].Text)
	}
	return "SOLICITAR_CONOCIMIENTO_HARDWARE"
}

func (a *Agent) Act(action string) {
	fmt.Printf("🚀 Decisión tomada por la IA: %s\n", action)

	if action == "SOLICITAR_CONOCIMIENTO_HARDWARE" {
		// El agente decide usar de forma autónoma la herramienta MCP
		doc := a.CallKnowledgeServer("Reglas de temperatura crítica GPU")
		fmt.Printf("📖 [Datos Obtenidos del Servidor]: %s\n", doc)
		// Guardamos la información recuperada en el estado del agente para el siguiente turno
		a.State.CurrentData = doc
	} else if action == "EJECUTAR_APAGADO_SEGURIDAD" {
		fmt.Println("🚨 ACCIÓN LOCAL: Mitigando riesgo térmico. Ejecutando desconexión forzada del Rig.")
		fmt.Println("🏁 Simulación de Interoperabilidad MCP completada con éxito.")
		a.State.TargetMet = true
	}
}

func main() {
	fmt.Println("🤖 Iniciando Agente con Interoperabilidad de Herramientas (Nivel 4)...")
	
	agent := Agent{
		State: AgentState{Step: 0, TargetMet: false},
	}

	for !agent.State.TargetMet && agent.State.Step < 3 {
		agent.Perceive()
		action := agent.Think()
		agent.Act(action)
		time.Sleep(2 * time.Second)
	}
}

