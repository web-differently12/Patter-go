import React, { useState } from 'react';

export interface Marker {
  timeMs: number;
  timeLabel: string;
  type: 'friction' | 'rag_low_confidence' | 'success';
  label: string;
  detail: string;
}

export interface BANTScore {
  budget: { score: number; detail: string };
  authority: { score: number; detail: string };
  need: { score: number; detail: string };
  timing: { score: number; detail: string };
}

export const VoiceFailuresPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'reasons' | 'inspector'>('inspector');
  const [selectedCallId, setSelectedCallId] = useState<string>('call_98234_failure');

  const markers: Marker[] = [
    { timeMs: 12000, timeLabel: '00:12', type: 'friction', label: 'Interruption / Frustration', detail: 'Client exprime de l\'impatience ("Laissez-moi parler!")' },
    { timeMs: 34000, timeLabel: '00:34', type: 'rag_low_confidence', label: 'RAG Confiance Faible (0.62)', detail: 'Question sur l\'API v2 - absence d\'info dans la base' },
    { timeMs: 68000, timeLabel: '01:08', type: 'friction', label: 'Transfert Demandé', detail: 'Client demande explicitement un opérateur humain' },
    { timeMs: 95000, timeLabel: '01:35', type: 'success', label: 'Accord / RDV Validé', detail: 'RDV confirmé pour mardi 14h' },
  ];

  const bantScore: BANTScore = {
    budget: { score: 85, detail: 'Budget validé (>50k€/an)' },
    authority: { score: 90, detail: 'VP Sales (Décideur principal)' },
    need: { score: 95, detail: 'Besoin urgent de remplacer la stack legacy' },
    timing: { score: 70, detail: 'Projet prévu pour Q3 2026' },
  };

  return (
    <div className="space-y-6 text-slate-100 p-6 bg-[#08090a] min-h-screen font-sans">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b border-slate-800 pb-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white">Diagnostic des Échecs & Post-Mortem des Appels</h1>
          <p className="text-sm text-slate-400">Analyse détaillée des échecs, des demandes de transfert et inspection mot-à-mot des appels.</p>
        </div>
        <div className="flex bg-slate-900 border border-slate-800 rounded-lg p-1 text-xs">
          <button
            onClick={() => setActiveTab('inspector')}
            className={`px-3 py-1.5 rounded-md font-medium transition-colors ${activeTab === 'inspector' ? 'bg-cyan-500/20 text-cyan-400' : 'text-slate-400 hover:text-white'}`}
          >
            Call Inspector
          </button>
          <button
            onClick={() => setActiveTab('reasons')}
            className={`px-3 py-1.5 rounded-md font-medium transition-colors ${activeTab === 'reasons' ? 'bg-cyan-500/20 text-cyan-400' : 'text-slate-400 hover:text-white'}`}
          >
            Top Raisons d'Échec
          </button>
        </div>
      </div>

      {activeTab === 'reasons' ? (
        /* Top Reasons Grid */
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
            <div className="flex items-center gap-2 mb-3">
              <span className="w-3 h-3 rounded-full bg-amber-400"></span>
              <h3 className="font-bold text-white text-base">Lacunes RAG / Connaissance</h3>
            </div>
            <p className="text-xs text-slate-400 mb-4">Questions posées sans réponse exacte dans la base de connaissances.</p>
            <div className="space-y-3">
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-amber-300 font-semibold mb-1">"Proposez-vous une conformité SOC2 Type II ?"</div>
                <div className="text-slate-400">42 appels concernés (Confiance RAG: 0.58)</div>
              </div>
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-amber-300 font-semibold mb-1">"Quels sont les tarifs pour +100k minutes/mois ?"</div>
                <div className="text-slate-400">28 appels concernés (Confiance RAG: 0.61)</div>
              </div>
            </div>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
            <div className="flex items-center gap-2 mb-3">
              <span className="w-3 h-3 rounded-full bg-rose-500"></span>
              <h3 className="font-bold text-white text-base">Interruption & Frustration</h3>
            </div>
            <p className="text-xs text-slate-400 mb-4">Mouvements d'irritation ou demande explicite d'intervenant humain.</p>
            <div className="space-y-3">
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-rose-300 font-semibold mb-1">"Passez-moi un conseiller humain immédiatement !"</div>
                <div className="text-slate-400">65 transferts forcés</div>
              </div>
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-rose-300 font-semibold mb-1">Interruption répétée du prompt IA</div>
                <div className="text-slate-400">31 appels écourtés</div>
              </div>
            </div>
          </div>

          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
            <div className="flex items-center gap-2 mb-3">
              <span className="w-3 h-3 rounded-full bg-slate-500"></span>
              <h3 className="font-bold text-white text-base">Répondeurs & Faux Numéros</h3>
            </div>
            <p className="text-xs text-slate-400 mb-4">Détection AMD (Answering Machine Detection) et non-joignabilité.</p>
            <div className="space-y-3">
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-slate-200 font-semibold mb-1">Messagerie vocale détectée</div>
                <div className="text-slate-400">184 appels (Voicemail drop effectué)</div>
              </div>
              <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs">
                <div className="text-slate-200 font-semibold mb-1">Numéro invalide / Non attribué</div>
                <div className="text-slate-400">12 appels rejetés par le carrier</div>
              </div>
            </div>
          </div>
        </div>
      ) : (
        /* Call Inspector */
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Player & Timeline */}
          <div className="lg:col-span-2 space-y-6">
            <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
              <div className="flex justify-between items-center mb-4">
                <span className="text-xs font-mono text-cyan-400 bg-cyan-950/60 px-2 py-1 rounded border border-cyan-800/50">{selectedCallId}</span>
                <span className="text-xs text-slate-400">Durée: 02:14 · Client: +33 6 12 34 56 78</span>
              </div>

              {/* Audio Waveform / Player Mock */}
              <div className="bg-slate-950 border border-slate-800 rounded-lg p-4 mb-4">
                <div className="flex items-center gap-4">
                  <button className="w-10 h-10 rounded-full bg-cyan-500 text-slate-950 font-bold flex items-center justify-center hover:bg-cyan-400">▶</button>
                  <div className="flex-1">
                    <div className="h-8 bg-slate-900 rounded flex items-center px-2 gap-1 relative overflow-hidden">
                      {/* Interactive Marker Dots */}
                      <span className="absolute left-[10%] w-3 h-3 rounded-full bg-rose-500 border border-white cursor-pointer" title="00:12 Friction"></span>
                      <span className="absolute left-[30%] w-3 h-3 rounded-full bg-amber-400 border border-white cursor-pointer" title="00:34 Confiance Faible"></span>
                      <span className="absolute left-[70%] w-3 h-3 rounded-full bg-rose-500 border border-white cursor-pointer" title="01:08 Transfert"></span>
                      <span className="absolute left-[90%] w-3 h-3 rounded-full bg-emerald-400 border border-white cursor-pointer" title="01:35 RDV Validé"></span>
                      <div className="w-1/3 bg-cyan-500/30 h-full"></div>
                    </div>
                  </div>
                  <span className="text-xs font-mono text-slate-400">00:45 / 02:14</span>
                </div>
              </div>

              {/* Visual Markers Timeline */}
              <h4 className="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-3">Marqueurs d'Inspection</h4>
              <div className="space-y-2">
                {markers.map((m, i) => (
                  <div key={i} className="flex items-start gap-3 p-2.5 rounded bg-slate-950/40 border border-slate-800/60 text-xs">
                    <span className={`px-2 py-0.5 rounded font-mono font-bold ${
                      m.type === 'friction' ? 'bg-rose-500/20 text-rose-400 border border-rose-500/30' :
                      m.type === 'rag_low_confidence' ? 'bg-amber-400/20 text-amber-300 border border-amber-400/30' :
                      'bg-emerald-400/20 text-emerald-400 border border-emerald-400/30'
                    }`}>
                      {m.timeLabel}
                    </span>
                    <div>
                      <div className="font-semibold text-slate-200">{m.label}</div>
                      <div className="text-slate-400 mt-0.5">{m.detail}</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Transcript Sync */}
            <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg">
              <h4 className="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-3">Transcription Synchronisée Mot-à-Mot</h4>
              <div className="space-y-3 text-xs max-h-60 overflow-y-auto pr-2">
                <div className="p-2 rounded bg-slate-950 border border-slate-800">
                  <span className="font-bold text-cyan-400">Agent IA:</span> Bonjour ! Je suis l'assistant EchoFlow. Comment puis-je vous aider aujourd'hui ?
                </div>
                <div className="p-2 rounded bg-rose-950/30 border border-rose-900/50 text-rose-200">
                  <span className="font-bold text-rose-400">Client:</span> Écoutez, je n'ai pas le temps, j'ai une question précise sur votre API v2 !
                </div>
                <div className="p-2 rounded bg-amber-950/30 border border-amber-900/50 text-amber-200">
                  <span className="font-bold text-amber-400">Agent IA (RAG 0.62):</span> L'API v2 offre plusieurs endpoints. Laissez-moi vérifier les détails dans la documentation...
                </div>
              </div>
            </div>
          </div>

          {/* BANT Score Card */}
          <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
            <div className="flex justify-between items-center border-b border-slate-800 pb-3">
              <h3 className="font-bold text-white text-sm">Score BANT LeMUR</h3>
              <span className="text-xs font-semibold px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">Qualifié A+</span>
            </div>

            <div className="space-y-4 text-xs">
              <div>
                <div className="flex justify-between font-medium mb-1">
                  <span className="text-slate-300">Budget</span>
                  <span className="text-emerald-400 font-bold">{bantScore.budget.score}%</span>
                </div>
                <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mb-1">
                  <div className="bg-emerald-400 h-full" style={{ width: `${bantScore.budget.score}%` }}></div>
                </div>
                <span className="text-slate-400 text-[11px]">{bantScore.budget.detail}</span>
              </div>

              <div>
                <div className="flex justify-between font-medium mb-1">
                  <span className="text-slate-300">Authority</span>
                  <span className="text-emerald-400 font-bold">{bantScore.authority.score}%</span>
                </div>
                <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mb-1">
                  <div className="bg-emerald-400 h-full" style={{ width: `${bantScore.authority.score}%` }}></div>
                </div>
                <span className="text-slate-400 text-[11px]">{bantScore.authority.detail}</span>
              </div>

              <div>
                <div className="flex justify-between font-medium mb-1">
                  <span className="text-slate-300">Need</span>
                  <span className="text-emerald-400 font-bold">{bantScore.need.score}%</span>
                </div>
                <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mb-1">
                  <div className="bg-emerald-400 h-full" style={{ width: `${bantScore.need.score}%` }}></div>
                </div>
                <span className="text-slate-400 text-[11px]">{bantScore.need.detail}</span>
              </div>

              <div>
                <div className="flex justify-between font-medium mb-1">
                  <span className="text-slate-300">Timing</span>
                  <span className="text-amber-400 font-bold">{bantScore.timing.score}%</span>
                </div>
                <div className="w-full bg-slate-800 h-1.5 rounded-full overflow-hidden mb-1">
                  <div className="bg-amber-400 h-full" style={{ width: `${bantScore.timing.score}%` }}></div>
                </div>
                <span className="text-slate-400 text-[11px]">{bantScore.timing.detail}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
