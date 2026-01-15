import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t bg-gray-50 dark:bg-gray-900">
      <div className="container py-8 md:py-12">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          <div>
            <h3 className="text-sm font-semibold mb-4">Платформ</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <Link href="/about" className="hover:text-foreground">
                  Бидний тухай
                </Link>
              </li>
              <li>
                <Link href="/features" className="hover:text-foreground">
                  Боломжууд
                </Link>
              </li>
              <li>
                <Link href="/pricing" className="hover:text-foreground">
                  Үнийн санал
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-4">Тусламж</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <Link href="/help" className="hover:text-foreground">
                  Тусламж
                </Link>
              </li>
              <li>
                <Link href="/docs" className="hover:text-foreground">
                  Гарын авлага
                </Link>
              </li>
              <li>
                <Link href="/faq" className="hover:text-foreground">
                  Түгээмэл асуултууд
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-4">Хууль эрх зүй</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>
                <Link href="/privacy" className="hover:text-foreground">
                  Нууцлалын бодлого
                </Link>
              </li>
              <li>
                <Link href="/terms" className="hover:text-foreground">
                  Үйлчилгээний нөхцөл
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h3 className="text-sm font-semibold mb-4">Холбоо барих</h3>
            <ul className="space-y-2 text-sm text-muted-foreground">
              <li>support@gerege.mn</li>
              <li>+976 7000-0000</li>
            </ul>
          </div>
        </div>
        <div className="mt-8 pt-8 border-t text-center text-sm text-muted-foreground">
          <p>&copy; 2024 Gebase Platform. Бүх эрх хуулиар хамгаалагдсан.</p>
        </div>
      </div>
    </footer>
  );
}
