import { Card, CardContent } from "@/components/ui/card";
import { Sprout, DollarSign, MapPin, Grid3X3, TrendingUp } from "lucide-react";

interface Stats {
  totalFarms: number;
  totalPlots: number;
  totalProduction: number;
  totalCost: number;
  coe?: number;
  cot?: number;
  ctProducao?: number;
}

export function StatsCards({ stats }: { stats: Stats }) {
  const items = [
    { icon: MapPin, label: "Fazendas", value: stats.totalFarms },
    { icon: Grid3X3, label: "Talhões", value: stats.totalPlots },
    {
      icon: Sprout,
      label: "Produção (sacas)",
      value: stats.totalProduction.toLocaleString(),
    },
    {
      icon: DollarSign,
      label: "Custo total",
      value: `R$ ${stats.totalCost.toLocaleString()}`,
    },
    ...(stats.coe
      ? [
          {
            icon: TrendingUp,
            label: "COE (Custo Operacional Efetivo)",
            value: `R$ ${stats.coe.toLocaleString()}`,
            title:
              "Custo Operacional Efetivo: soma dos desembolsos operacionais diretos (insumos, mão de obra contratada, combustível etc.)",
          },
        ]
      : []),
    ...(stats.cot
      ? [
          {
            icon: TrendingUp,
            label: "COT (Custo Operacional Total)",
            value: `R$ ${stats.cot.toLocaleString()}`,
            title:
              "Custo Operacional Total: COE + mão de obra familiar + depreciação de máquinas e benfeitorias",
          },
        ]
      : []),
    ...(stats.ctProducao
      ? [
          {
            icon: TrendingUp,
            label: "CT (Custo Total)",
            value: `R$ ${stats.ctProducao.toLocaleString()}`,
            title:
              "Custo Total: COT + remuneração do capital investido (terra e outros ativos)",
          },
        ]
      : []),
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {items.map((item) => (
        <Card key={item.label} title={item.title}>
          <CardContent className="flex items-center gap-4 p-4">
            <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
              <item.icon className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-muted-foreground">{item.label}</p>
              <p className="text-xl font-bold text-primary">{item.value}</p>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
