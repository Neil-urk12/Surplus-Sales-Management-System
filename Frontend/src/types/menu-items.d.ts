export interface menuItemsProps {
  title: string;
  icon?: string;
  to?: string;
  exact?: boolean;
  children?: menuItemsProps[];
  hasSubmenu?: boolean;
}
