import React, { useState } from 'react';

export const AgentEditorPage: React.FC = () => {
  // Conversational Prompts
  const [firstMessage, setFirstMessage] = useState('Bonjour ! Comment puis-je vous aider aujourd\'hui ?');
  const [systemPrompt, setSystemPrompt] = useState(
    'Vous etes l\'assistant accueil d\'entreprise. Repondez de maniere courtoise, synthetique et sans jargon.'
  );

  // Core Engines
  const [sttProvider, setSttProvider] = useState('Neural STT Ultra-Fast');
  const [model, setModel] = useState('Neural Realtime Engine');
  const [fallbackModel, setFallbackModel] = useState('Enterprise Fallback Model');
  const [voice, setVoice] = useState('Voice Flash Ultra-HD');
  const [fallbackVoice, setFallbackVoice] = useState('Voice Sonic HD (Fallback)');

  // Audio DSP & Turn Detection
  const [bargeInThresholdMs, setBargeInMs] = useState(300);
  const [backchanneling, setBackchanneling] = useState(true);
  const [ambientNoiseSound, setAmbientNoise] = useState<'office' | 'cafe' | 'silent'>('office');
  const [fillerWords, setFillerWords] = useState(true);
  const [noiseSuppression, setNoiseSuppression] = useState<'krisp' | 'rnnoise' | 'off'>('krisp');
  const [semanticTurnDetector, setSemanticTurnDetector] = useState<'smart_turn_v3' | 'namo' | 'vad_energy'>('smart_turn_v3');
  const [preemptiveGeneration, setPreemptiveGeneration] = useState(true);
  const [tokenCompaction, setTokenCompaction] = useState(true);

  // Webhooks & Server Events
  const [serverUrl, setServerUrl] = useState('https://api.patter.ai/v1/webhooks/voice');
  const [serverSecret, setServerSecret] = useState('sec_patter_live_9f8d721a');

  // Custom Tools / Functions Schema
  const [toolsSchema, setToolsSchema] = useState(
    JSON.stringify(
      [
        {
          name: 'check_calendar_availability',
          description: 'Verifier les creneaux disponibles pour un rendez-vous.',
          parameters: {
            type: 'object',
            properties: { date: { type: 'string', description: 'YYYY-MM-DD' } },
            required: ['date'],
          },
        },
      ],
      null,
      2
    )
  );

  // Structured Data Extraction Schema
  const [extractionSchema, setExtractionSchema] = useState(
    JSON.stringify(
      {
        customer_intent: 'booking | inquiry | complaint',
        lead_score: 'number (1-10)',
        next_action: 'string',
      },
      null,
      2
    )
  );

  // LockedInAI-inspired Realtime Visio Copilot Tuning
  const [copilotSuggestions, setCopilotSuggestions] = useState(true);
  const [objectionDetector, setObjectionDetector] = useState(true);

  // Lindy.ai-inspired Retail Workflow Triggers
  const [crmAutoSync, setCrmAutoSync] = useState(true);
  const [calendarAutoSchedule, setCalendarAutoSchedule] = useState(true);

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="border-b border-zinc-800 pb-4 flex justify-between items-center">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Agent Studio Configurator</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-0.5">Configuration Avancee de l'Agent Studio</h1>
        </div>
        <button className="bg-cyan-500 hover:bg-cyan-400 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          Enregistrer l'Agent
        </button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* System Prompt & Instructions */}
        <div className="lg:col-span-2 space-y-6">
          {/* Main Prompts */}
          <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4">
            <div>
              <label className="block text-xs font-semibold text-zinc-300 mb-1">Premier Message / Salutation Vocale (First Message)</label>
              <input
                type="text"
                value={firstMessage}
                onChange={(e) => setFirstMessage(e.target.value)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-xs text-zinc-200 font-mono focus:outline-none focus:border-cyan-500"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-zinc-300 mb-1">System Prompt & Directives Telephoniques</label>
              <textarea
                value={systemPrompt}
                onChange={(e) => setSystemPrompt(e.target.value)}
                rows={8}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-3 text-xs text-zinc-200 font-mono focus:outline-none focus:border-cyan-500 leading-relaxed"
              />
            </div>
          </div>

          {/* Custom Tools & Structured Data Extraction */}
          <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4">
            <h4 className="font-bold text-white text-xs">Custom Tools (Functions Calling) & Extraction Sémantique</h4>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs font-mono">
              <div className="space-y-1">
                <label className="text-zinc-400 block font-semibold">Outils Personnalises (Tools Schema JSON)</label>
                <textarea
                  value={toolsSchema}
                  onChange={(e) => setToolsSchema(e.target.value)}
                  rows={6}
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-[11px] text-cyan-300 focus:outline-none focus:border-cyan-500"
                />
              </div>

              <div className="space-y-1">
                <label className="text-zinc-400 block font-semibold">Schema d'Extraction Structuree Post-Appel</label>
                <textarea
                  value={extractionSchema}
                  onChange={(e) => setExtractionSchema(e.target.value)}
                  rows={6}
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-xl p-2.5 text-[11px] text-indigo-300 focus:outline-none focus:border-cyan-500"
                />
              </div>
            </div>
          </div>

          {/* Server Webhooks & Callbacks */}
          <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4 text-xs">
            <h4 className="font-bold text-white">Webhooks Serveur & Callback d'Evenements</h4>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 font-mono">
              <div>
                <label className="text-zinc-400 block mb-1">URL Webhook Serveur</label>
                <input
                  type="text"
                  value={serverUrl}
                  onChange={(e) => setServerUrl(e.target.value)}
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-lg p-2 text-zinc-200"
                />
              </div>

              <div>
                <label className="text-zinc-400 block mb-1">Secret de Signature Webhook</label>
                <input
                  type="password"
                  value={serverSecret}
                  onChange={(e) => setServerSecret(e.target.value)}
                  className="w-full bg-zinc-900 border border-zinc-800 rounded-lg p-2 text-zinc-200"
                />
              </div>
            </div>
          </div>

          {/* Audio DSP & Turn Detection Tuning */}
          <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-3 text-xs">
            <h4 className="font-bold text-white">Traitement du Signal Audio (DSP) & Detection de Parole</h4>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 font-mono">
              <div className="p-2 bg-zinc-900 rounded border border-zinc-800 space-y-1">
                <label className="text-[11px] text-zinc-400 block font-semibold">Suppression du Bruit</label>
                <select
                  value={noiseSuppression}
                  onChange={(e) => setNoiseSuppression(e.target.value as any)}
                  className="w-full bg-zinc-950 border border-zinc-700 rounded p-1 text-[11px] text-zinc-200"
                >
                  <option value="krisp">Krisp Neural AI (Haute Precision)</option>
                  <option value="rnnoise">RNNoise Open-Source (Leger)</option>
                  <option value="off">Desactive (Pass-Through Direct)</option>
                </select>
              </div>

              <div className="p-2 bg-zinc-900 rounded border border-zinc-800 space-y-1">
                <label className="text-[11px] text-zinc-400 block font-semibold">Detection Fin de Tour (Turn Detector)</label>
                <select
                  value={semanticTurnDetector}
                  onChange={(e) => setSemanticTurnDetector(e.target.value as any)}
                  className="w-full bg-zinc-950 border border-zinc-700 rounded p-1 text-[11px] text-zinc-200"
                >
                  <option value="smart_turn_v3">Smart-Turn v3 (Semantique ML)</option>
                  <option value="namo">NAMO Realtime Model</option>
                  <option value="vad_energy">VAD Energie Silences (Legacy)</option>
                </select>
              </div>

              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={preemptiveGeneration}
                  onChange={(e) => setPreemptiveGeneration(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Generation Preemptive LLM</span>
              </label>

              <label className="flex items-center gap-2 p-2 bg-zinc-900 rounded border border-zinc-800">
                <input
                  type="checkbox"
                  checked={tokenCompaction}
                  onChange={(e) => setTokenCompaction(e.target.checked)}
                  className="accent-cyan-500"
                />
                <span>Compaction de Tokens Contextuels</span>
              </label>
            </div>
          </div>
        </div>

        {/* Engine, Voice & Fallbacks */}
        <div className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 space-y-4 text-xs">
          <h3 className="font-bold text-white text-sm">Paramètres Moteurs, Voix & Fallback</h3>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Moteur Transcripteur STT</label>
            <select
              value={sttProvider}
              onChange={(e) => setSttProvider(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Neural STT Ultra-Fast">Neural STT Ultra-Fast (Sub-100ms)</option>
              <option value="Multilingual STT Enterprise">Multilingual STT Enterprise</option>
              <option value="On-Premise Whisper Engine">On-Premise Whisper Engine</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Modele LLM Principal</label>
            <select
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Neural Realtime Engine">Neural Realtime Engine</option>
              <option value="Enterprise Intelligence Tier">Enterprise Intelligence Tier</option>
              <option value="Realtime Sales Copilot Engine">Realtime Sales Copilot Engine</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">LLM de Fallback (Secours)</label>
            <select
              value={fallbackModel}
              onChange={(e) => setFallbackModel(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Enterprise Fallback Model">Enterprise Fallback Model</option>
              <option value="Fast Local Fallback">Fast Local Fallback Engine</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Voix Synthetique TTS</label>
            <select
              value={voice}
              onChange={(e) => setVoice(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Voice Flash Ultra-HD">Voice Flash Ultra-HD</option>
              <option value="Voice Sonic HD">Voice Sonic HD</option>
              <option value="Voice Alloy HD">Voice Alloy HD</option>
            </select>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Voix TTS de Fallback</label>
            <select
              value={fallbackVoice}
              onChange={(e) => setFallbackVoice(e.target.value)}
              className="w-full bg-zinc-900 border border-zinc-800 rounded-xl px-3 py-2 text-zinc-200"
            >
              <option value="Voice Sonic HD (Fallback)">Voice Sonic HD (Fallback)</option>
              <option value="Voice Standard Backup">Voice Standard Backup</option>
            </select>
          </div>

          {/* Fine Conversational Options */}
          <div className="border-t border-zinc-800/80 pt-3 space-y-2 font-mono">
            <span className="font-bold text-zinc-300 block">Fine Conversational Options</span>

            <label className="flex items-center gap-2 text-[11px]">
              <input
                type="checkbox"
                checked={backchanneling}
                onChange={(e) => setBackchanneling(e.target.checked)}
                className="accent-cyan-500"
              />
              <span>Backchanneling ("mm-hmm", "d'accord")</span>
            </label>

            <label className="flex items-center gap-2 text-[11px]">
              <input
                type="checkbox"
                checked={fillerWords}
                onChange={(e) => setFillerWords(e.target.checked)}
                className="accent-cyan-500"
              />
              <span>Mots de remplissage pendant RAG</span>
            </label>

            <div>
              <label className="block text-zinc-400 text-[11px] mb-1">Bruit d'ambiance d'arriere-plan</label>
              <select
                value={ambientNoiseSound}
                onChange={(e) => setAmbientNoise(e.target.value as any)}
                className="w-full bg-zinc-900 border border-zinc-800 rounded-lg px-2 py-1 text-zinc-300"
              >
                <option value="office">Ambiance Bureau Discret</option>
                <option value="cafe">Ambiance Cafe</option>
                <option value="silent">Silencieux (Insonorise)</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-zinc-400 font-medium mb-1">Seuil Interruption Barge-In: {bargeInThresholdMs}ms</label>
            <input
              type="range"
              min={100}
              max={1000}
              step={50}
              value={bargeInThresholdMs}
              onChange={(e) => setBargeInMs(Number(e.target.value))}
              className="w-full h-1.5 bg-zinc-800 rounded-lg appearance-none cursor-pointer accent-cyan-500"
            />
          </div>
        </div>
      </div>
    </div>
  );
};
