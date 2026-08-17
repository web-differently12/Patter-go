import React, { useState } from 'react';

export interface AgentItem {
  id: string;
  name: string;
  mode: 'realtime' | 'pipeline';
  voice: string;
  model: string;
  sttProvider?: string;
  ttsProvider?: string;
  status: 'active' | 'draft' | 'paused';
  totalTurns: number;
}

export const AgentsPage: React.FC = () => {
  const [agents] = useState<AgentItem[]>([
    {
      id: 'ag_receptionist_01',
      name: 'Agent Accueil Support Fr',
      mode: 'realtime',
      voice: 'alloy',
      model: 'gpt-realtime-mini',
      status: 'active',
      totalTurns: 12480,
    },
    {
      id: 'ag_qualifier_02',
      name: 'Agent Qualifier B2B SaaS',
      mode: 'pipeline',
      voice: 'eleven_flash_v2_5',
      model: 'claude-3-5-sonnet',
      sttProvider: 'Deepgram Nova-2',
      ttsProvider: 'ElevenLabs',
      status: 'active',
      totalTurns: 8420,
    },
    {
      id: 'ag_whatsapp_03',
      name: 'Assistant Messenger & WhatsApp',
      mode: 'pipeline',
      voice: 'cartesia_sonic',
      model: 'gpt-4o-mini',
      status: 'active',
      totalTurns: 19400,
    },
  ]);

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center border-b border-zinc-800 pb-4 gap-4">
        <div>
          <span className="text-xs font-mono font-bold text-cyan-400">Agent Studio & Fleet Management</span>
          <h1 className="text-2xl font-bold tracking-tight text-white mt-1">Gestion de la Flotte d'Agents IA</h1>
          <p className="text-xs text-slate-400">Concevez, testez et deployez vos agents vocaux, visio et messaging.</p>
        </div>

        <button className="bg-white hover:bg-zinc-200 text-zinc-950 font-bold px-4 py-2 rounded-xl text-xs transition-colors shadow-sm">
          + Creer un Nouvel Agent
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {agents.map((ag) => (
          <div key={ag.id} className="bg-zinc-950 border border-zinc-800 rounded-2xl p-5 shadow-xl flex flex-col justify-between space-y-4">
            <div>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-bold text-white text-base">{ag.name}</h3>
                  <span className="text-[10px] text-zinc-500 font-mono">{ag.id}</span>
                </div>
                <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-800">
                  {ag.status.toUpperCase()}
                </span>
              </div>

              <div className="space-y-1.5 text-xs text-slate-400 mt-4 font-mono">
                <div className="flex justify-between">
                  <span>Mode:</span>
                  <span className="text-cyan-400 font-bold">{ag.mode}</span>
                </div>
                <div className="flex justify-between">
                  <span>Modèle LLM:</span>
                  <span className="text-zinc-200">{ag.model}</span>
                </div>
                <div className="flex justify-between">
                  <span>Voix TTS:</span>
                  <span className="text-zinc-200">{ag.voice}</span>
                </div>
                {ag.sttProvider && (
                  <div className="flex justify-between">
                    <span>STT Provider:</span>
                    <span className="text-zinc-200">{ag.sttProvider}</span>
                  </div>
                )}
              </div>
            </div>

            <div className="pt-3 border-t border-zinc-800/80 flex justify-between items-center text-xs">
              <span className="text-zinc-500 font-mono">{ag.totalTurns.toLocaleString()} turns</span>
              <button className="text-cyan-400 hover:underline font-semibold">Éditer dans Agent Studio →</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
