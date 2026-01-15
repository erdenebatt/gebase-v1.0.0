"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import {
  HelpCircle,
  BookOpen,
  MessageCircle,
  Mail,
  Phone,
  ExternalLink,
  FileText,
  Video,
  Headphones,
  Clock,
} from "lucide-react";
import { Button } from "@/components/ui/button";

const faqs = [
  {
    question: "Хэрэглэгч хэрхэн нэмэх вэ?",
    answer:
      "Хэрэглэгчийн удирдлага -> Хэрэглэгчид хэсэгт орж, 'Шинэ хэрэглэгч' товч дарна. Шаардлагатай мэдээллийг бөглөж хадгална.",
  },
  {
    question: "Эрх хэрхэн тохируулах вэ?",
    answer:
      "Хэрэглэгчийн удирдлага -> Эрхүүд хэсэгт орж, эрхийн тохиргоог хийнэ. Эрх бүрт зөвшөөрөл болон цэсүүдийг холбож өгнө.",
  },
  {
    question: "Системүүдийг хэрхэн сэлгэх вэ?",
    answer:
      "Зүүн талын sidebar дээрх систем icon-ууд дээр дарж системүүдийг сэлгэнэ.",
  },
  {
    question: "Байгууллагын бүтэц хэрхэн тохируулах вэ?",
    answer:
      "Байгууллагын удирдлага -> Байгууллагууд хэсэгт орж, шаталсан бүтцийг үүсгэнэ. Эцэг байгууллага сонгож хүүхэд байгууллагуудыг нэмнэ.",
  },
  {
    question: "Цэсийг хэрхэн нэмэх вэ?",
    answer:
      "Системийн тохиргоо -> Цэсүүд хэсэгт орж шинэ цэс үүсгэнэ. Цэсийг эрхүүдтэй холбож харагдах байдлыг тохируулна.",
  },
  {
    question: "DSL схем гэж юу вэ?",
    answer:
      "DSL (Domain Specific Language) схем нь динамик өгөгдлийн бүтэц тодорхойлох боломж олгоно. Схем үүсгэж, талбаруудыг нэмж, дүрмүүд болон workflow тохируулна.",
  },
];

const resources = [
  {
    title: "Баримт бичиг",
    description: "Системийн бүрэн гарын авлага",
    icon: BookOpen,
    color: "text-blue-500",
    bgColor: "bg-blue-100 dark:bg-blue-900/30",
    action: "Гарын авлага үзэх",
  },
  {
    title: "Видео заавар",
    description: "Алхам алхмаар заавар",
    icon: Video,
    color: "text-purple-500",
    bgColor: "bg-purple-100 dark:bg-purple-900/30",
    action: "Видео үзэх",
  },
  {
    title: "Чат тусламж",
    description: "Шууд холбогдож асуулт асуух",
    icon: MessageCircle,
    color: "text-green-500",
    bgColor: "bg-green-100 dark:bg-green-900/30",
    action: "Чат эхлүүлэх",
  },
  {
    title: "Имэйл холбоо",
    description: "support@gerege.mn",
    icon: Mail,
    color: "text-amber-500",
    bgColor: "bg-amber-100 dark:bg-amber-900/30",
    action: "Имэйл илгээх",
  },
];

const stats = [
  {
    title: "Нийтлэл",
    value: "50+",
    icon: FileText,
  },
  {
    title: "Видео заавар",
    value: "25+",
    icon: Video,
  },
  {
    title: "Түгээмэл асуулт",
    value: `${faqs.length}`,
    icon: HelpCircle,
  },
  {
    title: "Дэмжлэг",
    value: "24/7",
    icon: Headphones,
  },
];

export default function HelpPage() {
  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
            <HelpCircle className="h-6 w-6 text-primary" />
            Тусламж
          </h1>
          <p className="text-muted-foreground">
            Gebase Platform-ийн хэрэглээний заавар болон тусламж
          </p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        {stats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Resources Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {resources.map((resource, index) => (
          <Card key={index} className="hover:shadow-md transition-shadow">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className={`p-3 rounded-lg ${resource.bgColor}`}>
                  <resource.icon className={`h-5 w-5 ${resource.color}`} />
                </div>
                <div>
                  <CardTitle className="text-base">{resource.title}</CardTitle>
                  <CardDescription className="text-xs">
                    {resource.description}
                  </CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <Button variant="outline" className="w-full gap-2">
                <ExternalLink className="h-4 w-4" />
                {resource.action}
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* FAQ Section */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <HelpCircle className="h-5 w-5" />
            </div>
            <div>
              <CardTitle>Түгээмэл асуултууд</CardTitle>
              <CardDescription>
                Хамгийн их асуудаг асуултууд болон хариултууд
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Accordion type="single" collapsible className="w-full">
            {faqs.map((faq, index) => (
              <AccordionItem key={index} value={`item-${index}`}>
                <AccordionTrigger className="text-left">
                  {faq.question}
                </AccordionTrigger>
                <AccordionContent className="text-muted-foreground">
                  {faq.answer}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </CardContent>
      </Card>

      {/* Contact Section */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900/30 text-blue-500">
                <Headphones className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Техникийн дэмжлэг</CardTitle>
                <CardDescription>
                  Техникийн асуудал шийдвэрлэх
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted">
              <Phone className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="font-medium">7777-1234</p>
                <p className="text-xs text-muted-foreground">Утасны дугаар</p>
              </div>
            </div>
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="font-medium">Даваа - Баасан</p>
                <p className="text-xs text-muted-foreground">09:00 - 18:00</p>
              </div>
            </div>
            <Button className="w-full gap-2">
              <Phone className="h-4 w-4" />
              Залгах
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 dark:bg-green-900/30 text-green-500">
                <Mail className="h-5 w-5" />
              </div>
              <div>
                <CardTitle>Борлуулалт</CardTitle>
                <CardDescription>
                  Захиалга болон санал хүсэлт
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted">
              <Phone className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="font-medium">7777-5678</p>
                <p className="text-xs text-muted-foreground">Утасны дугаар</p>
              </div>
            </div>
            <div className="flex items-center gap-3 p-3 rounded-lg bg-muted">
              <Mail className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="font-medium">sales@gerege.mn</p>
                <p className="text-xs text-muted-foreground">Имэйл хаяг</p>
              </div>
            </div>
            <Button variant="outline" className="w-full gap-2" asChild>
              <a href="mailto:sales@gerege.mn">
                <Mail className="h-4 w-4" />
                Имэйл илгээх
              </a>
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
